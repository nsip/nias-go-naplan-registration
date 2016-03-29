// aslvalidator reads in the asl data file
// construcutng checking keys based on known searches
// validation messages are checked to see if their ASL id is in the list
// and if the id is assigned to the expected state
//
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime"

	agg "github.com/nsip/nias-go-naplan-registration/aggregator/lib"
	"github.com/nats-io/nats"
	"github.com/wildducktheories/go-csv"
)

func main() {

	// handle command-line config options
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var vtype = flag.String("vtype", "ASL", "The type of validation, ie. content, business_logic, core etc.")
	var topic = flag.String("topic", "validation", "The root topic name to subscribe to")
	var qGroup = flag.String("qg", "aslvalidation", "The consumer group to join for parallel processing")
	var state = flag.String("state", "naplan", "The state identifier for this service [VIC, SA, NT, WA, ACT, TAS, NSW, QLD]")

	f, err := os.Open("./schoolslist/asl_schools.csv")
	reader := csv.WithIoReader(f)
	records, err := csv.ReadAll(reader)
	log.Printf("ASL records read: %v", len(records))
	if err != nil {
		log.Fatalf("Unable to open ASL schools file, service aborting...")
	}

	type searchterms struct {
		AcaraID string
		State   string
	}

	asl := make(map[string]searchterms)
	for _, r := range records {

		r := r.AsMap()

		s := searchterms{
			AcaraID: r["ACARA ID"],
			State:   r["State"],
		}
		asl[r["ACARA ID"]] = s
	}
	log.Println("...all ASL records ready for validation")

	// establish connection to NATS server
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("cannot reach NATS server, service will abort: ", err)
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	// listen on the subject channel for error messages
	_, err = nc.QueueSubscribe(*topic+"."+*state, *qGroup, func(msg *nats.Msg) {

		dat := make(map[string]string)
		if err := json.Unmarshal(msg.Data, &dat); err != nil {
			log.Println("Error unmarshalling json message: ", err)
		}
		// log.Println(dat)

		txID := dat["TxID"]

		// update on progress to monitors
		pn := agg.ProcessingNotification{txID, *vtype}
		ec.Publish("validation.status", pn)

		st, ok := asl[dat["ASLSchoolId"]]
		if !ok {
			desc := "ASL ID " + dat["ASLSchoolId"] + " not found in ASL list of valid IDs"
			msg := agg.ValidationError{
				Description:  desc,
				Field:        "ASLSchoolId",
				OriginalLine: dat["OriginalLine"],
				TxID:         txID,
				Vtype:        *vtype,
			}
			ec.Publish("validation.errors", msg)
		} else {
			if st.State != dat["StateTerritory"] {
				desc := "ASL ID " + dat["ASLSchoolId"] + " is as valid ID, but not for " + dat["StateTerritory"]
				msg := agg.ValidationError{
					Description:  desc,
					Field:        "ASLSchoolId",
					OriginalLine: dat["OriginalLine"],
					TxID:         txID,
					Vtype:        *vtype,
				}
				ec.Publish("validation.errors", msg)
			}
		}

	})
	if err != nil {
		log.Fatalf("cannot subscribe to required NATS queue, service will abort: %v\n", err)
	}

	log.Printf("Listening on [%s] as member of [%s]\n", *topic+"."+*state, *qGroup)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()

}
