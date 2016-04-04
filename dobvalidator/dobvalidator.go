// Date of birth validator
// Checks that provided birthdate for student is a valid date
// and that the dob, year level and test level all make sense
// according to the ACARA business rules as set out in the data spec
package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	agg "github.com/nsip/nias-go-naplan-registration/aggregator/lib"
        lib "github.com/nsip/nias-go-naplan-registration/lib"
	"github.com/nats-io/nats"
)

// check whether a date is within a range
func inDateRange(start, end, check time.Time) bool {
	// assume date range is inclusive
	if start.Equal(check) || end.Equal(check) {
		return true
	}
	return check.After(start) && check.Before(end)
}

// derive the year level from the given date of birth
func calculateYearLevel(t time.Time) string {

	yr := OOB
	if inDateRange(yr3start, yr3end, t) {
		// log.Println(t, "Matches Yr3: is between", yr3start, "and", yr3end, ".")
		return "3"
	}
	if inDateRange(yr5start, yr5end, t) {
		// log.Println(t, "Matches Yr3: is between", yr5start, "and", yr5end, ".")
		return "5"
	}
	if inDateRange(yr7start, yr7end, t) {
		// log.Println(t, "Matches Yr3: is between", yr7start, "and", yr7end, ".")
		return "7"
	}
	if inDateRange(yr9start, yr9end, t) {
		// log.Println(t, "Matches Yr3: is between", yr9start, "and", yr9end, ".")
		return "9"
	}
	// log.Println("Detected student year from dob is: ", yr)

	return yr

}

const OOB string = "out_of_band"

var yr3start, yr3end, yr5start, yr5end, yr7start, yr7end, yr9start, yr9end time.Time

func main() {
	// handle command-line config options
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display logging timestamps")
	var vtype = flag.String("vtype", "date", "The type of validation, ie. content, business_logic, core etc.")
	var topic = flag.String("topic", "validation", "The root topic name to subscribe to")
	var qGroup = flag.String("qg", "dobvalidation", "The consumer group to join for parallel processing")
	var state = flag.String("state", "naplan", "The state identifier for this service [VIC, SA, NT, WA, ACT, TAS, NSW, QLD]")
	var tstyr = flag.String("tstyr", "2017", "The year in which the test will occur; used to baseline the year/test level age range windows")

	log.SetFlags(0)
	flag.Parse()

	// warn if no flag set for tstyr

	// set up the calendar windows for testing against
	layout := "2006-01-02" // reference date format for date fields
	// create the testing 'now' date window baseline year
	tnow, _ := time.Parse(layout, *tstyr+"-01-01")
	yr3start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-9)+"-01-01")
	yr3end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-8)+"-07-31")

	yr5start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-11)+"-01-01")
	yr5end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-10)+"-07-31")

	yr7start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-13)+"-01-01")
	yr7end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-12)+"-07-31")

	yr9start, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-15)+"-01-01")
	yr9end, _ = time.Parse(layout, strconv.Itoa(tnow.Year()-14)+"-07-31")

	log.Println("\n==========================================\nNOTE:\n")
	log.Println("Baseline year for creating test/year level ranges is: ", *tstyr)
	log.Println("To change this, pass a year as comand line parameter -tstyr")
	log.Println("e.g. \n\n> dobvalidator(.exe) -tstyr 2018")
	log.Println("\n==========================================\n")
	log.Println("Date Validation Ranges:\n")
	log.Println("Year 3 range: ", yr3start.Format(layout), " - ", yr3end.Format(layout))
	log.Println("Year 5 range: ", yr5start.Format(layout), " - ", yr5end.Format(layout))
	log.Println("Year 7 range: ", yr7start.Format(layout), " - ", yr7end.Format(layout))
	log.Println("Year 9 range: ", yr9start.Format(layout), " - ", yr9end.Format(layout))
	log.Println("\n==========================================\n")

	// establish connection to NATS server
        natsconn := lib.NatsConn(*urls)

	// listen on the subject channel for messages & pass for validation
	_, err = natsconn.Nc.QueueSubscribe(*topic+"."+*state, *qGroup, func(msg *nats.Msg) {

		dat := make(map[string]string)
		if err := json.Unmarshal(msg.Data, &dat); err != nil {
			log.Println("Error unmarshalling json message: ", err)
		}

		// establish transaction id
		txID := dat["TxID"]

		// update status reporting
		pn := agg.ProcessingNotification{txID, *vtype}
		natsconn.Ec.Publish("validation.status", pn)

		t, err := time.Parse(layout, dat["BirthDate"])
		// log.Println("Provided birth date is: ", t)
		if err != nil {
			// log.Println("unable to parse date: ", err)
			desc := "Date provided does not parse correctly for yyyy-mm-dd"
			msg := agg.ValidationError{
				Description:  desc,
				Field:        "BirthDate",
				OriginalLine: dat["OriginalLine"],
				TxID:         txID,
				Vtype:        *vtype,
			}
			natsconn.Ec.Publish("validation.errors", msg)

		} else {

			yrlvl := dat["YearLevel"]
			desc := ""
			field := "BirthDate"
			ok := true
			switch {
			case yrlvl == "P":
				// log.Println("student is primary")
				desc = "Year level supplied is P, does not match expected test level " + dat["TestLevel"]
				field = field + "/TestLevel"
				ok = false
			case strings.Contains(yrlvl, "UG"):
				// log.Println("student is ungraded")
				desc = "Year level supplied is UG, will result in SRM warning flag for test level " + dat["TestLevel"]
				field = field + "/TestLevel/YearLevel"
				ok = false
			case yrlvl == "0":
				// log.Println("student is in year 0!!")
				desc = "Year level supplied is 0, does not match expected test level " + dat["TestLevel"]
				field = field + "/TestLevel"
				ok = false
			case calculateYearLevel(t) == OOB:
				// log.Println("Age derived year level not in any NAPLAN window")
				desc = "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges"
				field = field + "/YearLevel"
				ok = false
			default:
				field = "BirthDate"
				if yrlvl != calculateYearLevel(t) {
					// log.Println("Student is in wrong yr level: ", yrlvl)
					desc = "Student Year Level (yr " + yrlvl + ") does not match year level derived from BirthDate (yr " + calculateYearLevel(t) + ")"
					field = field + "/" + "YearLevel "
					ok = false
				}
				tstlvl, _ := dat["TestLevel"]
				if tstlvl != calculateYearLevel(t) {
					// log.Println("Student is in wrong test level: ", tstlvl)
					desc = "Student Test Level (yr " + tstlvl + ") does not match year level derived from BirthDate (yr " + calculateYearLevel(t) + ")"
					field = field + "/" + "TestLevel"
					ok = false
				}
			}

			if !ok {
				msg := agg.ValidationError{
					Description:  desc,
					Field:        field,
					OriginalLine: dat["OriginalLine"],
					TxID:         txID,
					Vtype:        *vtype,
				}
				natsconn.Ec.Publish("validation.errors", msg)
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
