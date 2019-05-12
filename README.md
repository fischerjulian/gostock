# GoStock

GoStock is a simple web services written in Go using the Iris web framework. The app is used to demonstrate
deployment to Cloud Foundry. Code quality and functionality is limited to the demo purpose.

## Deployment

How to deploy the app on Cloud Foundry ...

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