# Trafikverket-api golang client

## Usage

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"

    trv "github.com/roffe/trafikverket-api"
)

func main() {
    trv.Debug = true
    req := trv.NewRequest("your-apikey")
    req.Query(trv.Opts{
        "objecttype":    "TrainAnnouncement",
        "schemaversion": "1.6",
        "orderby":       "AdvertisedTimeAtLocation",
    },
        trv.Filter(
            trv.And(
                trv.Eq(trv.Opts{
                    "name":  "ActivityType",
                    "value": "Avgang"}),
                trv.Eq(trv.Opts{
                    "name":  "LocationSignature",
                    "value": "Cst"}),
                trv.Or(
                    trv.And(
                        trv.Gt(trv.Opts{
                            "name":  "AdvertisedTimeAtLocation",
                            "value": "$dateadd(-00:15:00)"},
                        ),
                        trv.Lt(trv.Opts{
                            "name":  "AdvertisedTimeAtLocation",
                            "value": "$dateadd(14:00:00)"},
                        ),
                    ),
                    trv.And(
                        trv.Lt(trv.Opts{
                            "name":  "AdvertisedTimeAtLocation",
                            "value": "$dateadd(00:30:00)"},
                        ),
                        trv.Gt(trv.Opts{
                            "name":  "EstimatedTimeAtLocation",
                            "value": "$dateadd(-00:15:00)"},
                        ),
                    ),
                ),
            ),
        ),
        trv.Include("AdvertisedTrainIdent"),
        trv.Include("AdvertisedTimeAtLocation"),
        trv.Include("TrackAtLocation"),
        trv.Include("ToLocation"),
    )

    resp, err := req.Do()
    if resp.Body != nil {
        defer resp.Body.Close()
    }
    if err != nil {
        log.Fatal(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(body[:]))
}
```
