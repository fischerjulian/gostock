# GoStock

GoStock is a simple web services written in Go using the Iris web framework. The app is used to demonstrate
deployment to Cloud Foundry. Code quality and functionality is limited to the demo purpose.

## Deployment

How to deploy the app on Cloud Foundry ...

### The Go Buildpack

This application makes use of the [go buildpack](https://docs.cloudfoundry.org/buildpacks/go/index.html). 
The automatic detection of the go buildpack works if the application uses [dep](https://github.com/golang/dep) for dependency management.

    cf logs gostock --recent

### Seeding

The application automatically creates two entries upon its first start.

## Usage

Once the application is running its API can be used as follows. Examples are provided as CURL commands.

### Retrieve Stock Entries

In order to retrieve stock entries as JSON execute the following CURL command:

    curl http://localhost:8080/stocks

### Create a new Stock Entry

In order to create a stock entry execute the following CURL command:

     curl -d „name=BASF&value=6523" -X POST http://localhost:8080/stock

Note that the value of stocks is given as EUR cents and therefore are represented as integers.