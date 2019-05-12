# GoStock

GoStock is a simple web services written in Go using the Iris web framework. The app is used to demonstrate
deployment to Cloud Foundry. Code quality and functionality is limited to the demo purpose.

## Deployment

How to deploy the app on Cloud Foundry platform is described in a separate article.
You might have to adapt application > name in manifest.yml in case the name has already been taken.

The application work on [paas.anynines](paas.anynines) using [a9s PostgreSQL](https://www.anynines.com/) as well as on run.pivotal.io using ElephantSQL.

### The Go Buildpack

This application makes use of the [go buildpack](https://docs.cloudfoundry.org/buildpacks/go/index.html). 
The automatic detection of the go buildpack works if the application uses [dep](https://github.com/golang/dep) for dependency management.

To workaround a problem with dep and Iris the following lines have been added to Gopkg.toml:

    [[override]]
    name = "github.com/flosch/pongo2"
    branch = "master"

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

## Local Development

### VCAP_SERVICES Credentials

If you want to develop locally you need to set environment VCAP variables such as:

    export VCAP_APPLICATION='{"application_id":"f142621f-9307-4c28-a3aa-fb263abbbab1","application_name":"gostock","application_uris":["gostock.de.a9sapp.eu"],"application_version":"6ae902e7-6639-44cf-bb9a-715a8eaccf93","cf_api":"https://api.de.a9s.eu","limits":{"disk":512,"fds":16384,"mem":256},"name":"gostock","space_id":"dbbdec3c-8191-4c82-a3f9-1fdd31cc5fbf","space_name":"production","uris":["gostock.de.a9sapp.eu"],"users":null,"version":"6ae902e7-6639-44cf-bb9a-715a8eaccf93“}'d

and

    export VCAP_SERVICES='{"a9s-postgresql10":[{"binding_name":null,"credentials":{"host":"localhost","hosts":["localhostc"],"name":"gostock","password":"","port":5432,"uri":"postgres://jfischer:@localhost:5432/gostock","username":"jfischer"},"instance_name":"gostockdb","label":"a9s-postgresql10","name":"gostock","plan":"postgresql-single-nano","provider":null,"syslog_drain_url":null,"tags":["sql","database","object-relational","consistent"],"volume_mounts":[]}]}'

The application will otherwise crash at startup.