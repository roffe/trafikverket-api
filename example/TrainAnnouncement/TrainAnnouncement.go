package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	trv "github.com/roffe/trafikverket-api"
)

var apiKey = "openapiconsolekey"

func main() {
	trv.Debug = true
	req := trv.NewRequest(
		apiKey,
		trv.Query(
			trv.Opts{
				trv.OptObjtype:       "TrainAnnouncement",
				trv.OptSchemaversion: "1.6",
				trv.OptOrderBy:       "AdvertisedTimeAtLocation",
			},
			trv.Filter(
				trv.And(
					trv.Eq(trv.Opts{
						trv.OptName:  "ActivityType",
						trv.OptValue: "Avgang"}),
					trv.Eq(trv.Opts{
						trv.OptName:  "LocationSignature",
						trv.OptValue: "Cst"}),
					trv.Or(
						trv.And(
							trv.Gt(trv.Opts{
								trv.OptName:  "AdvertisedTimeAtLocation",
								trv.OptValue: "$dateadd(-00:15:00)"},
							),
							trv.Lt(trv.Opts{
								trv.OptName:  "AdvertisedTimeAtLocation",
								trv.OptValue: "$dateadd(14:00:00)"},
							),
						),
						trv.And(
							trv.Lt(trv.Opts{
								trv.OptName:  "AdvertisedTimeAtLocation",
								trv.OptValue: "$dateadd(00:30:00)"},
							),
							trv.Gt(trv.Opts{
								trv.OptName:  "EstimatedTimeAtLocation",
								trv.OptValue: "$dateadd(-00:15:00)"},
							),
						),
					),
				),
			),
			trv.Include("AdvertisedTrainIdent"),
			trv.Include("AdvertisedTimeAtLocation"),
			trv.Include("TrackAtLocation"),
			trv.Include("ToLocation"),
		),
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

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatal(err)
	}
}
