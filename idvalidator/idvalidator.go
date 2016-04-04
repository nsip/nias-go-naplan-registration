// idvalidator.go
// service listens for messages on the NATS channel supplied as a
// startup command-line param
// all messages are passed for validation against created maps
// of checking structures
package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"

	agg "github.com/nsip/nias-go-naplan-registration/aggregator/lib"
        lib "github.com/nsip/nias-go-naplan-registration/lib"
	"github.com/nats-io/nats"
)

func main() {
	// handle command-line config options
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var vtype = flag.String("vtype", "identity", "The type of validation, ie. content, business_logic, core etc.")
	var topic = flag.String("topic", "validation", "The root topic name to subscribe to")
	var qGroup = flag.String("qg", "idvalidation", "The consumer group to join for parallel processing")
	var state = flag.String("state", "naplan", "The state identifier for this service [VIC, SA, NT, WA, ACT, TAS, NSW, QLD]")
	// var jsonSchema = flag.String("schema", "core.json", "The schema file to be used for validation by this instance of the validator")

	log.SetFlags(0)
	flag.Parse()

	// establish connection to NATS server
        natsconn := lib.NatsConn(*urls)

	// set up the data structures to be used
	// simple duplicate check is, have we seen this userid for this school before
	type simplekey struct {
		LocalId     string
		ASLSchoolId string
	}

	// this checks the user against a set of likely colliding matches
	type extendedkey struct {
		LocalId     string
		ASLSchoolId string
		FamilyName  string
		GivenName   string
		BirthDate   string
	}

	// data stores have outer key of transaction-id
	dsimple := make(map[string]map[simplekey]string)
	dcomplex := make(map[string]map[extendedkey]string)

	// listen on the subject channel for messages & pass for validation
	_, err := natsconn.Nc.QueueSubscribe(*topic+"."+*state, *qGroup, func(msg *nats.Msg) {

		dat := make(map[string]string)
		if err := json.Unmarshal(msg.Data, &dat); err != nil {
			log.Println("Error unmarshalling json message: ", err)
		}
		// log.Println(dat)

		txID := dat["TxID"]

		pn := agg.ProcessingNotification{txID, *vtype}
		natsconn.Ec.Publish("validation.status", pn)

		k := simplekey{
			LocalId:     dat["LocalId"],
			ASLSchoolId: dat["ASLSchoolId"],
		}

		// log.Println(k)

		// if the combination for simplekey not seen before then add it,
		// if it has then raise potential duplicate error
		ol, ok := dsimple[txID][k]
		if !ok {
			// log.Println("added new record")
			_, ok := dsimple[txID]
			if !ok {
				m := make(map[simplekey]string)
				dsimple[txID] = m
			}
			dsimple[txID][k] = dat["OriginalLine"]
		} else {
			// log.Println("duplicate of record: ", ol)
			desc := "LocalID (Student) and ASL ID (School) are potential duplicate of record: " + ol
			msg := agg.ValidationError{
				Description:  desc,
				Field:        "LocalID/ASL ID",
				OriginalLine: dat["OriginalLine"],
				TxID:         txID,
				Vtype:        *vtype,
			}
			natsconn.Ec.Publish("validation.errors", msg)
		}

		// more complex match
		ek := extendedkey{
			LocalId:     dat["LocalId"],
			ASLSchoolId: dat["ASLSchoolId"],
			FamilyName:  dat["FamilyName"],
			GivenName:   dat["GivenName"],
			BirthDate:   dat["BirthDate"],
		}

		// log.Println(k)

		// if the combination for complexkey not seen before then add it,
		// if it has then raise potential duplicate error
		ol, ok = dcomplex[txID][ek]
		if !ok {
			// log.Println("added new record")
			_, ok := dcomplex[txID]
			if !ok {
				m := make(map[extendedkey]string)
				dcomplex[txID] = m
			}
			dcomplex[txID][ek] = dat["OriginalLine"]
		} else {
			// log.Println("duplicate of record: ", ol)
			desc := "Potential duplicate of record: " + ol + "\n" +
				"based on matching: student local id, school asl id, family & given names and birthdate"
			msg := agg.ValidationError{
				Description:  desc,
				Field:        "Multiple (see description)",
				OriginalLine: dat["OriginalLine"],
				TxID:         txID,
				Vtype:        *vtype,
			}
			natsconn.Ec.Publish("validation.errors", msg)
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
