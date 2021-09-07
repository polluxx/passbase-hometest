You need to do some steps in order to run the solution:
- `go mod download`
- `cd cmd/rates-sync & go build`
- `./rates-sync run` to run the rates sync (you can add there some flags also)
- `cd ../server & go build`
- `./server` to run a server
- enjoy :)

In order to get your API token, you need to register a project:

Example payload to register a project

     curl -X POST \
    http://localhost:8080/project \
    -H 'content-type: application/json' \
    -d '{ "name": "Test project", "customer_email": "test@t.t" }'

Example payload to get conversion results

    curl http://localhost:8080/convert\?token\=YOUR-TOKEN\&source=usd\&destination=eur\&amount=120

OR you can try our swagger API :):
[check this out](http://localhost:8080/swagger/index.html)

To run tests:
    
    go test -v ./...

