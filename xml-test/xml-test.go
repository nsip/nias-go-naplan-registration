package main

import (
	"fmt"
	"flag"
	"os"
	xml "github.com/nsip/nias-go-naplan-registration/xml"
)

var inputFile = flag.String("infile", "5000students.xml", "Input file path")

func main() {
	fmt.Println("Hello world")

	flag.Parse()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	xml.XmlParse(xmlFile)
}

