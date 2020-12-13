# CoralReportApi
Runs the queries against the bigquery and saves the reports which can be served through Rest Apis

## Where it can be used

If you are building an app or page for reporting data. You save time on building the complete infrastructure for fetching and saving the data.
Lot of big data systems are not design for synchronous apis. So you have to build systems to retrieve the big data and cache it locally for faster and repeated access.


### Getting started
- create google project and a service account. downloads credentials and do `export GOOGLE_APPLICATION_CREDENTIALS=<CREDENTIALS_FILE>`
- clone git package
- Modify `resources/application-dev.properties` with project id
- run `docker build . -t coralreportapi:latest` 
- run `docker run -it --rm --name coralreportapi -p 8080:8080 -v $GOOGLE_APPLICATION_CREDENTIALS:/tmp/creds.json:ro -e GOOGLE_APPLICATION_CREDENTIALS=/tmp/creds.json coralreportapi:latest`

### Contact us
