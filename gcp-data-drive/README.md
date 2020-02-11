# gcp-data-drive
GCP Data Drive leverages a simple composable url path to retrieve data in JSON formats from supported GCP data platforms. bigquery and firestore are currently supported.

## Installation
### Deploy to appengine

:warning: Edit the deploy script and app.yaml below for your needs. Used as it, the command below will deploy and new version of the default appengine service.
```bash
git clone https://github.com/GoogleCloudPlatform/DIY-Tools.git
cd DIY-Tools/gcp-data-drive
./appengine_deploy.sh {gcp_projectid}

```

## web api composition
Each web api is composed by a simple drive navigation pattern.
https://{host}/{platform}/{gcp_project}/{param1}/param2}...

## Examples

### Bigquery
Assuming projectid of testbqproject in bigquery dataset mybqviews with a view name of coolnumbersview cloud be accessed via the following gcp-data-drive path:
https://{host}/bq/testbqproject/mybqviews/collnumbersview

### Firestore
Assuming projectid of testfsproject in firestore collection firstcollection in document firstdocument with collection mydocs all documents would be returned with following gcp-data-drive path:
https://{host}/fs/testfsproject/firstcollection/firstdocument/mydocs

A single document can also be accessed with the following:
https://{host}/fs/testfsproject/firstcollection/firstdocument/mydocs/12345

## Authentication
When deployed on appengine, the app engine default service account must be granted bigquery read and bigquery create job permission. These settings are the default if the appengine service and firestore or bigquery are in the same project.
