// simple csv-sif xml conversion routine
// receives a csv file, and returns content as SIF StudentPersonal xml

package main

import (
	"flag"
	"github.com/wildducktheories/go-csv"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/twinj/uuid"
)

// xml template is compiled once at startup and stored in sptmpl for reuse
var sptmpl *template.Template
var maxRecords int

func main() {

	flag.IntVar(&maxRecords, "maxr", 50000, "The number of records that will be processed in a single file")
	log.SetFlags(0)
	flag.Parse()

	log.Println("Initialising uuid generator")
	config := uuid.StateSaverConfig{SaveReport: true, SaveSchedule: 30 * time.Minute}
	uuid.SetupFileSystemStateSaver(config)
	log.Println("UUID generator initialised.")

	log.Println("Loading xml templates")
	fp := path.Join("templates", "studentpersonals.tmpl")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatalf("Unable to parse xml template, service aborting...")
	}
	sptmpl = tmpl
	log.Println("Template loaded ok.")
	log.Println("Server up; listening on :3000/convert")

	assets := http.StripPrefix("/", http.FileServer(http.Dir("public/")))

	http.Handle("/", assets)
	http.HandleFunc("/convert", convert)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// handler typically called from a web form frontend that sends files of any size using
// multipart form, will detect ajax style calls also by looking in header for file attachment,
// if multipart/form handles as uploaded file, if just data returns data stream
func convert(w http.ResponseWriter, r *http.Request) {

	var reader csv.Reader
	var fname, rfname string

	isMultipart := true
	file, hdr, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrNotMultipart {
			isMultipart = false
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

	if isMultipart {
		reader = csv.WithIoReader(file)
		fname = hdr.Filename
		rplcr := strings.NewReplacer(".csv", ".xml")
		rfname = rplcr.Replace(fname)
	} else {
		reader = csv.WithIoReader(r.Body)
	}

	// read the csv file
	records, err := csv.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Records received: ", len(records))

	if len(records) > maxRecords {
		http.Error(w, "Too many records submitted.\nTo stay within SRM file upload limits, maximum no. students in a single file is "+strconv.Itoa(maxRecords), http.StatusBadRequest)
		return
	}

	// create valid guids
	spsnls := make([]map[string]string, 0)
	for _, r := range records {
		r := r.AsMap()
		r["SIFuuid"] = uuid.NewV4().String()
		spsnls = append(spsnls, r)
	}

	// set headers to 'force' file download where appropriate
	if isMultipart {
		w.Header().Set("Content-Disposition", "attachment; filename="+rfname)
	}
	w.Header().Set("Content-Type", "application/xml")

	// apply the template & write to the client
	if err := sptmpl.Execute(w, spsnls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
