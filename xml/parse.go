// Shared code - XML Parser

package xml

import (
	"fmt"
	"io"
	"os"
	"encoding/xml"
	"encoding/json"
)

var OUTPUT = false;

// XXX add elements
// XXX Transaction ID and Sequence ID (see CSV)
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

				records = append(records, p);

				total++
			}
		default:
		}
		
	}

	fmt.Printf("Total articles: %d \n", total)

	return(records)
}

