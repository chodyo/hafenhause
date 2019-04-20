# Hafenhause

An experimental project where I learn how to use gcloud tools and create some neat applications for my household while I'm at it.

## Features

### Bedtimes

A simple bedtime tracker/display. Good behavior one night earns 10m later while bad behavior earns 10m earlier bedtime the next night.

## Server

### Things learned
- Go modules
- gcloud functions, firestore
- basic Go HTTP

### Deploying & testing

To deploy a new update to a function, run from the server directory

`bin/deploy.sh [functionName]`

(for example, `bin/deploy.sh Bedtime`)

To test, open [hafenhause_tests.http](server/bin/hafenhause_tests.http) and use the REST Client VSCode extension to send POST requests.

### Functions

#### Bedtime

1. CREATE: `/Bedtime/name`

Creates a person with the default bedtime.
   
2. READ: `/Bedtime/[name]`

Gets the requested person's bedtimes, or all bedtimes if no name is requested.

``` json
// response
[{
    "name": "Cody",
    "hour": 23,
    "minute": 59
},{
    "name": "Julia",
    "hour": 19,
    "minute": 30
}]
```

3. UPDATE: `/Bedtime/name`

Updates the person's bedtime, or all bedtimes if no name is specified.

``` json
// request
{
    "hour": 20,
    "minute": 0
}
```
   
4. DELETE: `/Bedtime/name`

Deletes the person's bedtime.
