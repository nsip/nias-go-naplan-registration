package main

import (
	"fmt"
	"flag"
	"os"
	xml "github.com/nsip/nias-go-naplan-registration/xml"
)

var inputFile = flag.String("infile", "5000students.xml", "Input file path")

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
		// r := removeBlanks(r.AsMap())
		fmt.Println(i);
		fmt.Println(r.GivenName);
	}
}

