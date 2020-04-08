# gcp-data-drive
GCP Data Drive leverages a composable url path to retrieve data in JSON formats from supported GCP data platforms. Bigquery and Firestore are currently supported.

## Installation
Clone this repo. Change directory to gcp-data-drive.
```bash
git clone https://github.com/GoogleCloudPlatform/DIY-Tools
cd DIY-Tools/gcp-data-drive
```

### Deploy to multiple compute platforms  
[Google Cloud SDK](https://cloud.google.com/sdk/docs/quickstarts) must be installed and initialized to preform these actions. Edit the PROJECT_ID variable as required.

#### Cloud Run
```bash
PROJECT_ID=YourProjectID
gcloud builds submit --config cloudbuild_run.yaml \
--project $PROJECT_ID --no-source \
--substitutions=_GIT_SOURCE_BRANCH="master",_GIT_SOURCE_URL="https://github.com/GoogleCloudPlatform/DIY-Tools"
```

#### Cloud Functions
```bash
PROJECT_ID=YourProjectID
gcloud builds submit --config cloudbuild_gcf.yaml \
--project $PROJECT_ID --no-source \
--substitutions=_GIT_SOURCE_BRANCH="master",_GIT_SOURCE_URL="https://github.com/GoogleCloudPlatform/DIY-Tools"
```


#### Appengine
:warning: Edit the deploy script and app.yaml for your needs. Used as is, the command below will deploy and new version of the default App Engine service. Fork this repo to customize your app.yaml configuration.

```bash
PROJECT_ID=YourProjectID
gcloud builds submit  --config cloudbuild_appengine.yaml \
--project $PROJECT_ID --no-source \
--substitutions=_GIT_SOURCE_BRANCH="master",_GIT_SOURCE_URL="https://github.com/GoogleCloudPlatform/DIY-Tools"
```

## Web API Composition
Each web api is composed by a drive navigation pattern.
https://{host}/{platform}/{gcp_project}/{param1}/param2}...

## Examples

### Bigquery
Assuming projectid of testbqproject in Bigquery dataset mybqviews with a view name of coolnumbersview cloud be accessed via the following gcp-data-drive path:
https://{host}/bq/testbqproject/mybqviews/collnumbersview

### Firestore
Assuming projectid of testfsproject in Firestore collection firstcollection in document firstdocument with collection mydocs all documents would be returned with following gcp-data-drive path:
https://{host}/fs/testfsproject/firstcollection/firstdocument/mydocs

A single document can also be accessed with the following:
https://{host}/fs/testfsproject/firstcollection/firstdocument/mydocs/12345

## Authentication
When deployed on App Engine, the app engine default service account must be granted Bigquery read and Bigquery create job permission. These settings are the default if the App Engine service and Firestore or Bigquery are in the same project.
