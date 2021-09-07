package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"math"
	"net/http"
	_ "passbase-hometest/cmd/server/docs"
	"passbase-hometest/domain"
	"passbase-hometest/domain/database"
	"strconv"

	"go.uber.org/zap"
)

var logger = zap.S().Named("router")

type Router struct {
	Engine            *gin.Engine
	repositoryService database.Repository
}

// @title Passbase Home test API
// @version 1.0
// @description This is a simple REST server for currency conversions.

// @contact.name Yuriy Kosakivsky
// @contact.url https://www.eliftech.com
// @contact.email yuriy.kosakivsky@eliftech.com

// @BasePath /
func NewRouter(repositoryService database.Repository) *Router {
	router := gin.Default()

	gr := &Router{
		Engine:            router,
		repositoryService: repositoryService,
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")

	router.POST("/project", gr.createProject)

	authorized := router.Group("/")
	authorized.Use(AuthMiddleware(func(token string) bool {
		project, err := repositoryService.FindProjectByToken(context.Background(), token)
		if project == nil || err != nil {
			logger.Errorf("token %s was not found; err: %s", token, err)
			return false
		}
		return true
	}))
	{
		authorized.GET("/convert", gr.convertRates)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return gr
}

// @Description signing up a new project and returns an unique api token
// @Accept  json
// @Produce  json
// @Param   Project.Name     body    string     true        "Project name"
// @Param   Project.CustomerEmail     body    string     true        "Unique customer email"
// @Success 200 {string} json	"{"token": "your token"}"
// @Failure 500 {object} string "internal server error"
// @Failure 400 {object} string "project with that email is already signed up"
// @Router /project [post]
func (r *Router) createProject(c *gin.Context) {
	var payload Project
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := r.repositoryService.FindProjectByEmail(context.Background(), payload.CustomerEmail)
	if err != nil {
		logger.With("email", payload.CustomerEmail).Errorf("error on project validation: %s", err)
		fmt.Printf("error on project validation: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if project != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project with that email is already signed up"})
		return
	}

	newToken, _ := GenerateRandomStringURLSafe(32)
	domainProject := domain.Project{
		Name:          payload.Name,
		CustomerEmail: payload.CustomerEmail,
		Token:         newToken,
	}
	_, err = r.repositoryService.RegisterProject(context.Background(), domainProject)
	if err != nil {
		logger.With("email", payload.CustomerEmail).Errorf("error on project creation: %s", err)
		fmt.Printf("error on project creation: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error happen on project creation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are signed up", "token": newToken})
}

// @Description converts passed amount based on source and destination currency rates
// @Accept  json
// @Produce  json
// @Param   token     query    string     true        "your API token"
// @Param   source     query    string     true        "source currency"
// @Param   destination      query    string     true        "destination currency"
// @Param   amount      query    int     true        "source amount"
// @Success 200 {string} json	"{"converted": {"currency": "destination currency", "amount": "amount"}}"
// @Failure 500 {object} string "internal server error"
// @Failure 400 {object} string "project with that email is already signed up"
// @Router /convert [get]
func (r *Router) convertRates(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")
	amount := c.Query("amount")

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("can't convert amount param: %s", err)})
		return
	}

	sourceCurrency := domain.CurrencyToDomain(source)
	if sourceCurrency == domain.RateUndefined {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("currency is not defined: %s", source)})
		return
	}

	dstCurrency := domain.CurrencyToDomain(destination)
	if dstCurrency == domain.RateUndefined {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("currency is not defined: %s", destination)})
		return
	}

	baseCurrency, err := r.repositoryService.GetBaseCurrency(context.Background())
	if err != nil {
		logger.Errorf("error on get base currency: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	sourceCurr, err := r.repositoryService.GetCurrencyRate(context.Background(), sourceCurrency.ToLabel())
	if err != nil {
		logger.Errorf("error on get source currency: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	destCurr, err := r.repositoryService.GetCurrencyRate(context.Background(), dstCurrency.ToLabel())
	if err != nil {
		logger.Errorf("error on get destination currency: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	convertedSum := math.Floor(sourceCurr.Convert(*destCurr, amountInt))

	c.JSON(http.StatusOK,
		gin.H{
			"base": baseCurrency.Currency,
			"converted": gin.H{
				"amount":   convertedSum,
				"currency": destCurr.Currency,
			}})
}
