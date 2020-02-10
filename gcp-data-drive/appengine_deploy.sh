cd cmd/webserver
gcloud app deploy app.yaml --project $1 -q
cd -