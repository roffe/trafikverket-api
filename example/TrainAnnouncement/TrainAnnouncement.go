package main

import (
	"encoding/json"
	"log"

	trv "github.com/roffe/trafikverket-api"
)

var apiKey = "openapiconsolekey"

func main() {
	trv.Debug = true
	req := trv.NewRequest(
		apiKey,
		trv.Query().Opts(trv.Opts{
			trv.OptObjtype:       "TrainAnnouncement",
			trv.OptSchemaversion: "1.6",
			trv.OptOrderBy:       "AdvertisedTimeAtLocation",
		}).Tags(
			trv.Filter().Tags(
				trv.And().Tags(
					trv.Eq().Opts(trv.Opts{
						trv.OptName:  "ActivityType",
						trv.OptValue: "Avgang"}),
					trv.Eq().Opts(trv.Opts{
						trv.OptName:  "LocationSignature",
						trv.OptValue: "Cst"}),
					trv.Or().Tags(
						trv.And().Tags(
							trv.Gt().Opts(trv.Opts{
								trv.OptName:  "AdvertisedTimeAtLocation",
								trv.OptValue: "$dateadd(-00:15:00)"},
							),
							trv.Lt().Opts(trv.Opts{
								trv.OptName:  "AdvertisedTimeAtLocation",
								trv.OptValue: "$dateadd(14:00:00)"},
							),
						),
						trv.And().Tags(
							trv.Lt().Opts(trv.Opts{
								trv.OptName:  "AdvertisedTimeAtLocation",
								trv.OptValue: "$dateadd(00:30:00)"},
							),
							trv.Gt().Opts(trv.Opts{
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
	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resp, &res); err != nil {
		log.Fatal(err)
	}
}
