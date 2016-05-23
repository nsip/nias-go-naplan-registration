// Shared code - XML Parser

package xml

import (
	"fmt"
	"os"
	"encoding/xml"
	"encoding/json"
)

var OUTPUT = true;

// XXX add elements
// XXX Transaction ID and Sequence ID (see CSV)
/*

CSV Files Headers

LocalId,SectorId,DiocesanId,OtherId,TAAId,StateProvinceId,NationalId,PlatformId,PreviousLocalId,PreviousSectorId,PreviousDiocesanId,PreviousOtherId,PreviousTAAId,PreviousStateProvinceId,PreviousNationalId,PreviousPlatformId,FamilyName,GivenName,PreferredName,MiddleName,BirthDate,Sex,CountryOfBirth,EducationSupport,FFPOS,VisaCode,IndigenousStatus,LBOTE,StudentLOTE,YearLevel,TestLevel,FTE,Homegroup,ClassCode,ASLSchoolId,SchoolLocalId,LocalCampusId,MainSchoolFlag,OtherSchoolId,ReportingSchoolId,HomeSchooledStudent,Sensitive,OfflineDelivery,Parent1SchoolEducation,Parent1NonSchoolEducation,Parent1Occupation,Parent1LOTE,Parent2SchoolEducation,Parent2NonSchoolEducation,Parent2Occupation,Parent2LOTE,AddressLine1,AddressLine2,Locality,Postcode,StateTerritory

*/
type student struct {
	XMLName   xml.Name `xml:"StudentPersonal"`
	LocalId string `xml:"LocalId" json:"ID"`
	StateProvinceId string `xml:"StateProvinceId"`
	FamilyName string `xml:"PersonInfo>Name>FamilyName"`
	GivenName string `xml:"PersonInfo>Name>GivenName"`
	MiddleName string `xml:"PersonInfo>Name>MiddleName"`
	PreferredGivenName string `xml:"PersonInfo>Name>PreferredGivenName"`
	HappyBirthDay string `xml:"PersonInfo>Demographics>BirthDate"`
}

// Function call
//	- Data file (XML)
//	- Type (Struct ?)
//	- Queue (for adding to nats)
//	- automatic detect
func XmlParse(xmlFile *os.File) {
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

				if (OUTPUT) {
					fmt.Println(
						"ID =", p.LocalId, 
						", Name =", p.GivenName, p.FamilyName, 
						"BirthDate=", p.HappyBirthDay)

					fmt.Print("\nXML=");
					x_enc.Encode(p)
					fmt.Print("\nJSON=");
					j_enc.Encode(p)
				}

				total++
			}
		default:
		}
		
	}

	fmt.Printf("Total articles: %d \n", total)
}

