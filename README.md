# Hafenhaus

An experimental project where I learn how to use gcloud tools and create some neat applications for my household while I'm at it.

## Features

### Bedtimes

A simple bedtime tracker/display. Good behavior one night earns 10m later while bad behavior earns 10m earlier bedtime the next night.

## Server

### Things learned
- Go modules
- gcloud functions, firestore

### Deploying & testing

To deploy a new update to a function, run from the server directory

`bin/deploy.sh [functionName]`

(for example, `bin/deploy.sh SubmitBedtimeReport`)

To test, open [hafenhaus_tests.http](server/bin/hafenhaus_tests.http) and use the REST Client VSCode extension to send POST requests.

### Functions

1. `GET /GetBedtimes`

2. `POST /SubmitBedtimeReport`

``` json
// SubmitBedtimeReport Request

// Set the time explicitly
{
    "subject": "cody",
    "date": "2019-04-07T19:30:00-04:00"
}
// Record the score and let the server figure out the new time
{
    "subject": "cody",
    "score": 10
}
```
