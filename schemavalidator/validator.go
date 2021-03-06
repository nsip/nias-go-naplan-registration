// validator.go
// service listens for messages on the NATS channel supplied as a
// startup command-line param (has sensible default if not provided)
// all messages are passed for validation using a json schema file
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/kardianos/osext"
	"github.com/nats-io/nats"
	agg "github.com/nsip/nias-go-naplan-registration/aggregator/lib"
	"github.com/xeipuuv/gojsonschema"
)

func main() {
	// handle command-line config options
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var vtype = flag.String("vtype", "content", "The type of validation, ie. content, business_logic, core etc.")
	var topic = flag.String("topic", "validation", "The root topic name to subscribe to")
	var qGroup = flag.String("qg", "csvvalidation", "The consumer group to join for parallel processing")
	var state = flag.String("state", "naplan", "The state identifier for this service [VIC, SA, NT, WA, ACT, TAS, NSW, QLD]")
	var jsonSchema = flag.String("schema", "core.json", "The schema file to be used for validation by this instance of the validator")

	exeDir, _ := osext.ExecutableFolder()
	log.Println(exeDir)

	log.SetFlags(0)
	flag.Parse()

	// establish connection to NATS server
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("cannot reach NATS server, service will abort: %v\n", err)
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	// load the validation schema once for efficiency
	loadSchema := func() *gojsonschema.Schema {

		s, readerr := ioutil.ReadFile("schemas/" + *jsonSchema)
		if readerr != nil {
			log.Fatalf("Unable to open schema file, service aborting...")
		}
		schemaLoader := gojsonschema.NewStringLoader(string(s))
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil {
			panic("schema load error in startup - \n\n" + err.Error() + "\n...service aborting.")
		}
		log.Println("loaded schema - schemas/" + *jsonSchema)
		return schema
	}
	schema := loadSchema()

	// listen on the subject channel for messages & pass for validation
	_, err = nc.QueueSubscribe(*topic+"."+*state, *qGroup, func(msg *nats.Msg) {

		dat := make(map[string]string)
		if err := json.Unmarshal(msg.Data, &dat); err != nil {
			log.Println("Error unmarshalling json message: ", err)
		}
		// log.Println(dat)

		pn := agg.ProcessingNotification{dat["TxID"], *vtype}
		ec.Publish("validation.status", pn)

		payloadLoader := gojsonschema.NewStringLoader(string(msg.Data))

		result, err := schema.Validate(payloadLoader)
		if err != nil {
			log.Println("Message validation schema processing error: " + err.Error())
		}

		if !result.Valid() {

			for _, desc := range result.Errors() {

				// trap enum errors for large enums such as country code
				// and truncate the message to prevent unwieldy message
				if desc.Type() == "enum" {
					switch desc.Field() {
					case "CountryOfBirth":
						desc.SetDescription("Country Code is not one of SACC 1269.0 codeset")
					case "Parent1LOTE", "Parent2LOTE", "StudentLOTE":
						desc.SetDescription("Language Code is not one of ASCL 1267.0 codeset")
					case "VisaCode":
						desc.SetDescription("Visa Code is not one of known values from http://www.immi.gov.au")
					}

				}

				msg := agg.ValidationError{
					Description:  desc.Description(),
					Field:        desc.Field(),
					OriginalLine: dat["OriginalLine"],
					TxID:         dat["TxID"],
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
