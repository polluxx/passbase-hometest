package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"passbase-hometest/config"
	"passbase-hometest/domain/database"
	"passbase-hometest/domain/providers"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	configFolder   = "../../config"
	updateInterval = int64(12 * 60)
	configuration  *config.Config
	repository     database.Repository
	ratesProvider  providers.Providerer

	log = zap.S().Named("rates sync")
)

func main() {
	app := cli.NewApp()
	app.Name = "rates-sync"

	app.Usage = `This app is used to sync the DB information with the rates results from 3d-party service.`

	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Folder where config.yml is located",
			//TakesFile:   true,
			Value:       configFolder,
			Destination: &configFolder,
		},
		&cli.Int64Flag{
			Name:        "interval",
			Usage:       "Update interval for between sync calls; sets in minutes",
			Value:       updateInterval,
			Destination: &updateInterval,
		},
	}

	app.Flags = commonFlags

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "Start rates sync process",
			Flags: commonFlags,
			Action: func(ctx *cli.Context) error {
				setup()
				syncWithTimer()
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(fmt.Sprintf("can't run sync app: %s", err))
	}
}

func setup() {
	var err error
	flag.Parse()

	if updateInterval < 1 {
		log.Fatalf("update interval %s can't be less than 1 minute", updateInterval)
	}

	configuration, err = config.LoadFromPath(configFolder)
	if err != nil {
		log.Fatalf("can't load config: %s", err)
	}

	if configuration == nil {
		log.Fatal("config is nil")
	}

	// we may use any rates provider if the implementation fits the provider interface,
	// by that we guarantee backwards compatibility and flexibility to switch to another one
	ratesProvider, err = providers.New(providers.Fixer, configuration)
	if err != nil {
		log.Fatalf("can't init the rates provider: %s", err)
	}

	repository = database.Connect(configuration.Database, database.SQLite)
}

func syncWithTimer() {
	syncTime := time.Duration(updateInterval) * time.Second
	timer := time.NewTimer(syncTime)
	// call this first time before timer will run
	err := sync()
	if err != nil {
		log.Fatalf("can't get rates: %s", err)
	}
	for {
		timer.Reset(syncTime)
		if <-timer.C; true {
			err = sync()
			if err != nil {
				log.Fatalf("can't get rates: %s", err)
			}
		}
	}
}

func sync() error {
	results, err := ratesProvider.Latest()
	if err != nil || len(results) == 0 {
		log.Errorf("can't get the latest rates from provider: %s, or results are empty", err)
		return err
	}

	// close db connections if 1 minute passed
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	for _, rate := range results {
		dbRate, errDB := repository.GetCurrencyRate(ctx, rate.Currency)
		errRecordNotFound := errors.Is(errDB, database.ErrRecordNotFound)
		if errDB != nil && !errRecordNotFound {
			log.Errorf("can't get the rate from DB: %s", errDB)
			return errDB
		}

		fmt.Println("rate", dbRate.Currency, ":", dbRate.Value)

		if dbRate == nil || errRecordNotFound {
			_, errDB = repository.CreateRate(ctx, rate)
			if errDB != nil {
				log.Errorf("can't insert the rate into DB: %s", errDB)
				return errDB
			}
			// we created a new record - continue to the next iteration
			continue
		}

		if !dbRate.Updated.Before(rate.Updated) {
			fmt.Printf("skip updating currency rate %s: up to date: %s\n", rate.Currency, rate.Updated)
			continue
		}

		// if record is found and no errors occurs - updating the record
		errDB = repository.UpdateRate(ctx, rate)
		if errDB != nil {
			log.Errorf("can't update the rate in DB: %s", errDB)
			return errDB
		}
	}

	fmt.Printf("Synced! Next iteration in: %s \n", time.Duration(updateInterval)*time.Second)
	return nil
}
