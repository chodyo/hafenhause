# hafenhaus

An experimental project where I learn how to use gcloud tools and create some neat applications for my household while I'm at it.

## things learned
- Go modules
- gcloud functions, firestore

## deploying, testing

To deploy a new update, run

`bin/deploy.sh [functionName]`
(for example, `bin/deploy.sh SubmitBedtimeReport`)

To test, open `hafenhaus_tests.http` and use the REST Client VSCode extension to send POST requests.

## functions
1. `SubmitBedtimeReport`

Request
``` json
{
	"subject": "cody",
	"date": "2019-04-07T19:10:00-04:00",
	"score": 10,
	"carryOver": true
}
```

## bedtimes

A bedtime tracker. Good behavior one night earns 10m later while bad behavior earns 10m earlier bedtime the next night.

### todo
| type  | description                                                                 |   status   |
| :---: | :-------------------------------------------------------------------------- | :--------: |
|  DB   | hold member states                                                          | **[DONE]** |
|  FN   | can upload new bedtime data each night                                      | **[DONE]** |
|  FN   | Background watcher/cron or additional functionality to update member states | **[TODO]** |
|  FN   | get current bedtime for each member                                         | **[TODO]** |
