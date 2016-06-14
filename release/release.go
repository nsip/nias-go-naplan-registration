package main

import (
	"encoding/json"
	"gopkg.in/cheggaaa/pb.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Release Information
type jsonRelease struct {
	Id        int    `json:"id"`
	TagName   string `json:"tag_name"`
	UploadURL string `json:"upload_url"`
}

type jsonUpload struct {
	DownloadURL string `json:"browser_download_url"`
}

// XXX Variables / Config
//	Base URL
//	Project Name
// 	How to get next release name/number

func getRelease(project string) jsonRelease {
	release := jsonRelease{}

	// Get the current Releases
	req, err := http.NewRequest("GET", "https://api.github.com/repos/nsip/"+project+"/releases/latest", nil)
	if err != nil {
		// handle err
	}
	// XXX DO NOT RELEASE !!!
	req.SetBasicAuth("scottp", "XXX")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// log.Printf("Body: %s\n", body)

	err = json.Unmarshal(body, &release)
	if err != nil {
		log.Fatal(err)
	}

	return release
}

func uploadFile(release jsonRelease, name string, filename string) jsonUpload {
	upload := jsonUpload{}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	uploadURLs := strings.Split(release.UploadURL, "{")
	uploadURL := uploadURLs[0]
	log.Printf("Sending file to %s", uploadURL+"?name="+name)

	bar := pb.StartNew(int(fi.Size())).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	pr := bar.NewProxyReader(f)

	req, err := http.NewRequest("POST", uploadURL+"?name="+name, pr)
	if err != nil {
		// handle err
		log.Printf("Error upload = %s", err)
	}
	req.SetBasicAuth("scottp", "XXX")
	req.Header.Set("Content-Type", "application/zip")
	req.ContentLength = fi.Size()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		// handle err
		log.Printf("Error upload = %s", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	bar.Finish()
	err = json.Unmarshal(body, &upload)
	if err != nil {
		log.Fatal(err)
	}

	if len(upload.DownloadURL) > 0 {
		log.Printf("Download URL = %s", upload.DownloadURL)
	} else {
		log.Printf("Body: %s\n", body)
	}
	return upload
}

// Arguments
//	1 = Project name
//	2 = File name to upload
//	3 = Local file path
func main() {
	// Get Latest Release
	release := getRelease(os.Args[1])

	log.Printf("Received release %d as %s", release.Id, release.TagName)

	uploadFile(release, os.Args[2], os.Args[3])
}