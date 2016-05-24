// Shared code - XML Parser

package xml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

var OUTPUT = false

// XXX add elements
// XXX Transaction ID and Sequence ID (see CSV)
/*

CSV Files Headers

LocalId,SectorId,DiocesanId,OtherId,TAAId,StateProvinceId,NationalId,PlatformId,PreviousLocalId,PreviousSectorId,PreviousDiocesanId,PreviousOtherId,PreviousTAAId,PreviousStateProvinceId,PreviousNationalId,PreviousPlatformId,FamilyName,GivenName,PreferredName,MiddleName,BirthDate,Sex,CountryOfBirth,EducationSupport,FFPOS,VisaCode,IndigenousStatus,LBOTE,StudentLOTE,YearLevel,TestLevel,FTE,Homegroup,ClassCode,ASLSchoolId,SchoolLocalId,LocalCampusId,MainSchoolFlag,OtherSchoolId,ReportingSchoolId,HomeSchooledStudent,Sensitive,OfflineDelivery,Parent1SchoolEducation,Parent1NonSchoolEducation,Parent1Occupation,Parent1LOTE,Parent2SchoolEducation,Parent2NonSchoolEducation,Parent2Occupation,Parent2LOTE,AddressLine1,AddressLine2,Locality,Postcode,StateTerritory

*/
type student struct {
	// XML Configuration
	// XMLName            xml.Name `xml:"StudentPersonal"`

	// Internal data
	TxID         string
	OriginalLine string

	// Important fields
	ASLSchoolId               string `json:",omitempty"`
	AddressLine1              string `json:",omitempty" xml:"PersonInfo>AddressList>Street>Address>Line1"`
	AddressLine2              string `json:",omitempty" xml:"PersonInfo>AddressList>Street>Address>Line2"`
	BirthDate                 string `json:",omitempty" xml:"PersonInfo>Demographics>BirthDate"`
	ClassCode                 string `json:",omitempty" xml:"MostRecent/ClassCode"`
	CountryOfBirth            string `json:",omitempty" xml:"PersonInfo>Demographics>CountryOfBirth"`
	DiocesanId                string `json:",omitempty"`
	EducationSupport          string `json:",omitempty"`
	FFPOS                     string `json:",omitempty" xml:"MostRecent>FFPOS"`
	FTE                       string `json:",omitempty" xml:"MostRecent>FTE"`
	FamilyName                string `json:",omitempty" xml:"PersonInfo>Name>FamilyName"`
	GivenName                 string `json:",omitempty" xml:"PersonInfo>Name>GivenName"`
	HomeSchooledStudent       string `json:",omitempty"`
	Homegroup                 string `json:",omitempty"`
	IndigenousStatus          string `json:",omitempty"`
	JurisdictionId            string `json:",omitempty"`
	LBOTE                     string `json:",omitempty" xml:"PersonInfo>Demographics>LBOTE"`
	LocalCampusId             string `json:",omitempty"`
	LocalId                   string `json:",omitempty" xml:"LocalId" json:"ID"`
	Locality                  string `json:",omitempty"`
	MainSchoolFlag            string `json:",omitempty"`
	MiddleName                string `json:",omitempty" xml:"PersonInfo>Name>MiddleName"`
	NationalId                string `json:",omitempty"`
	OfflineDelivery           string `json:",omitempty"`
	OtherId                   string `json:",omitempty"`
	OtherSchoolId             string `json:",omitempty"`
	Parent1LOTE               string `json:",omitempty" xml:"MostRecent>Parent1Language"` // Wrong?
	Parent1NonSchoolEducation string `json:",omitempty" xml:"MostRecent>Parent1EmploymentType"`
	Parent1Occupation         string `json:",omitempty"`
	Parent1SchoolEducation    string `json:",omitempty"`
	Parent2LOTE               string `json:",omitempty" xml:"MostRecent>Parent2Language"`
	Parent2NonSchoolEducation string `json:",omitempty"`
	Parent2Occupation         string `json:",omitempty" xml:"MostRecent>Parent2EmploymentType"`
	Parent2SchoolEducation    string `json:",omitempty"`
	PlatformId                string `json:",omitempty"`
	Postcode                  string `json:",omitempty" xml:"PersonInfo>AddressList>Address>PostalCode`
	PreferredName             string `json:",omitempty" xml:"PersonInfo>Name>PreferredGivenName"`
	PreviousDiocesanId        string `json:",omitempty"`
	PreviousJurisdictionId    string `json:",omitempty"`
	PreviousLocalId           string `json:",omitempty"`
	PreviousNationalId        string `json:",omitempty"`
	PreviousOtherId           string `json:",omitempty"`
	PreviousPlatformId        string `json:",omitempty"`
	PreviousSectorId          string `json:",omitempty"`
	PreviousStateProvinceId   string `json:",omitempty"`
	PreviousTAAId             string `json:",omitempty"`
	ReportingSchoolId         string `json:",omitempty"`
	SchoolLocalId             string `json:",omitempty"`
	SectorId                  string `json:",omitempty"`
	Sensitive                 string `json:",omitempty"`
	Sex                       string `json:",omitempty" xml:"PersonInfo>Demographics>Sex"`
	StateProvinceId           string `json:",omitempty" xml:"StateProvinceId"`
	StateTerritory            string `json:",omitempty"`
	StudentLOTE               string `json:",omitempty"`
	TAAId                     string `json:",omitempty"`
	TestLevel                 string `json:",omitempty" xml:"MostRecent>TestLevel>Code"`
	VisaCode                  string `json:",omitempty" xml:"PersonInfo>Demographics>VisaSubClass"`
	YearLevel                 string `json:",omitempty" xml:"MostRecent>YearLevel>Code"`
}

// Function call
//	- Data file (XML)
//	- Type (Struct ?)
//	- Queue (for adding to nats)
//	- automatic detect
func XmlParse(xmlFile io.Reader) (records []student) {
	// XML Encoder
	x_enc := xml.NewEncoder(os.Stdout)
	x_enc.Indent("  ", "    ")

	// JSON Encoder
	j_enc := json.NewEncoder(os.Stdout)
	// j_enc.Indent("  ", "    ")

	decoder := xml.NewDecoder(xmlFile)
	total := 0
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "StudentPersonal" {
				var p student
				decoder.DecodeElement(&p, &se)

				if OUTPUT {
					fmt.Println(
						"ID =", p.LocalId,
						", Name =", p.GivenName, p.FamilyName,
						"BirthDate=", p.BirthDate)
					fmt.Print("\nXML=")
					x_enc.Encode(p)
					fmt.Print("\nJSON=")
					j_enc.Encode(p)
				}

				records = append(records, p)

				total++
			}
		default:
		}

	}

	fmt.Printf("Total articles: %d \n", total)

	return (records)
}
