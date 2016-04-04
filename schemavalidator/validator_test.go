package main

import (
	"os"
	"io/ioutil"
	"testing"
	"bytes"
        "strconv"
	"time"
        "encoding/json"
	"log"
	"strings"

	"github.com/nats-io/nats"
	lib "github.com/nsip/nias-go-naplan-registration/lib"
        "github.com/wildducktheories/go-csv"
)

var natsconn *lib.NatsConnection

func csv2nats(csvstring string) ([]map[string]string, error) {
	// input csv, and output a slice of records to push to NATS. the same way that aggregator Post("/naplan/reg/:stateID") does
	reader := csv.WithIoReader(ioutil.NopCloser(bytes.NewReader([]byte(csvstring))))
	records, err := csv.ReadAll(reader)
	ret := make([]map[string]string, len(records))
	for i, r := range records {
		ret[i] = r.AsMap()
		ret[i]["OriginalLine"] = strconv.Itoa(i + 1)
		ret[i]["TxID"] = "dummyTxID"
	}
	return ret, err
}

func TestSex(t *testing.T) {

	csv := `LocalId,SectorId,DiocesanId,OtherId,TAAId,JurisdictionId,NationalId,PlatformId,PreviousLocalId,PreviousSectorId,PreviousDiocesanId,PreviousOtherId,PreviousTAAId,PreviousJurisdictionId,PreviousNationalId,PreviousPlatformId,FamilyName,GivenName,PreferredName,MiddleName,BirthDate,Sex,CountryOfBirth,EducationSupport,FFPOS,VisaCode,IndigenousStatus,LBOTE,StudentLOTE,YearLevel,TestLevel,FTE,HomeGroup,ClassCode,ASLSchoolId,SchoolLocalId,LocalCampusId,MainSchoolFlag,OtherSchoolId,ReportingSchoolId,HomeSchooledStudent,Sensitive,OfflineDelivery,Parent1SchoolEducation,Parent1NonSchoolEducation,Parent1Occupation,Parent1LOTE,Parent2SchoolEducation,Parent2NonSchoolEducation,Parent2Occupation,Parent2LOTE,AddressLine1,AddressLine2,Locality,Postcode,StateTerritory
fjghh425,14668,65616,75189,50668,59286,35164,47618,66065,4716,50001,65241,55578,44128,37734,73143,Seefeldt,Treva,Treva,E,2004-07-26,6,1101,Y,1,101,2,Y,2201,7,7,0.89,7E,7D,48096,046129,01,02,48096,48096,U,Y,Y,3,8,2,1201,2,7,4,1201,30769 PineTree Rd.,,Pepper Pike,9999,QLD`

	records, err := csv2nats(csv)
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	log.Printf("records received: %v", len(records))
	sub, err := natsconn.Nc.SubscribeSync("validation.errors") 
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	for _, r := range records {
		log.Println(r)
		natsconn.Ec.Publish("validation.naplan", r)
	}
	msg, err := sub.NextMsg(3 * time.Second) 
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	dat := make(map[string]string)
	if err := json.Unmarshal(msg.Data, &dat); err != nil {
		t.Fatalf("Error unmarshalling json message: %s", err)
	}
	log.Println(dat)
	if dat["errField"] != "Sex" {
		t.Fatalf("Expected error field %s, got field %s", "Sex", dat["errField"])
	}
	if !strings.Contains(dat["description"], "Sex must be one of the following") {
		t.Fatalf("Expected error description %s, got description %s", "Sex must be one of the following...", dat["description"])
	}
}

func TestMain(m *testing.M) {
	natsconn = lib.NatsConn(nats.DefaultURL)
	os.Exit(m.Run())
}
