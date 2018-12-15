# Reddit snapshots
A small GoLang program used to take snapshots of the current state of Reddit subreddits. Currently meant to be ran in AWS Lambda or as a cron job in general, with plans to make it runnable in a container with a perpetual timer or something of the sort.

Run with the following environment variables for configuration:

- CLIENT_ID
- CLIENT_SECRET
- USERNAME
- PASSWORD
- DB_URL
- DB_NAME (The name of the database)
- SNAPSHOTS_COLLECTION (Tje collection in which the snapshots will be stored)
- CONFIG_COLLECTION (The collection in which the configuration will be stored)

## Example for configuration:
```
{ 
    "_id" : ObjectId("5c122eb7a1783f58255167bd"), 
    "subreddits" : [
        {
            "subreddit" : "dankmemes", 
            "sort" : "top"
        }, 
        {
            "subreddit" : "worldnews", 
            "sort" : "hot"
        }, 
        {
            "subreddit" : "space", 
            "sort" : "hot"
        }, 
    ]
}
```
