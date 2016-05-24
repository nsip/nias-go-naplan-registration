package main

import (
	"flag"
	"fmt"
	xml "github.com/nsip/nias-go-naplan-registration/xml"
	"os"
)

var inputFile = flag.String("infile", "test_data/5000students.xml", "Input file path")

func main() {
	fmt.Println("Hello world - testing XML Parser")

	flag.Parse()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	records := xml.XmlParse(xmlFile)
	for i, r := range records {
		fmt.Print(i)
		// fmt.Println(r.FamilyName)
		fmt.Printf("%+v\n", r)
	}
}
