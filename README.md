# hafenhaus

An experimental project where I learn how to use gcloud tools and create some neat applications for my household while I'm at it.

## deploying, testing

To deploy a new update, run

`bin/deploy.sh [functionName]`
(for example, `bin/deploy.sh SubmitBedtimeReport`)

To test, open `hafenhaus_tests.http` and use the REST Client VSCode extension to send POST requests.

## bedtimes

A bedtime tracker. Good behavior one night earns 10m later while bad behavior earns 10m earlier bedtime the next night.

### todo
| type  | description                                                                 |      status       |
| :---: | :-------------------------------------------------------------------------- | :---------------: |
|  DB   | hold bedtime entries                                                        |    **[DONE]**     |
|  FN   | upload new bedtime data each night                                          | **[IN PROGRESS]** |
|  DB   | hold member states                                                          |    **[TODO]**     |
|  FN   | Background watcher/cron or additional functionality to update member states |    **[TODO]**     |
|  FN   | get current bedtime for each member                                         |    **[TODO]**     |
|  FN   | edit or delete data                                                         |    **[TODO]**     |

1. upload new bedtime data
   - separate the document out by members
   - possibly skip straight to holding state
   - don't round timestamp to nearest day - just accept events (maybe throttle to prevent spamming or multiple sequential requests)
   - keep in mind this is a document store, not a sql table

### things learned
- Go modules
- gcloud functions, firestore
