# GoStock

GoStock is a simple web services written in Go using the Iris web framework. The app is used to demonstrate
deployment to Cloud Foundry. Code quality and functionality is limited to the demo purpose.

## Deployment

How to deploy the app on Cloud Foundry ...

## Usage

Once the application is running its API can be used as follows. Examples are provided as CURL commands.

     curl -d „stock=BASF&value=6523" -X POST http://localhost:8000/stock
