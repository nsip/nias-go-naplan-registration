// web service to provide html/javascript ui for
// the various validation services
// also acts as reverse proxy to web services to keep
// endpoint management simple for clients
package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var port = ":8080"

var server_configs = map[string]string{
	"converterURL":      "http://localhost:3000", // optional, provide credentials through local.ini
	"aggregatorURL":     "http://localhost:1324", // optional
	"static_web_assets": "public/",               // static files
	"http_port":         port,
}

func main() {

	flag.StringVar(&port, "port", ":8080", "The port to run this webservice on.")
	log.SetFlags(0)
	flag.Parse()
	server_configs["http_port"] = port

	mux := http.NewServeMux()
	// by default, URL's will be mapped to our static assets
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(server_configs["static_web_assets"]))))
	mux.Handle("/nias/", http.StripPrefix("/nias/", http.FileServer(http.Dir(server_configs["static_web_assets"]))))

	// create a reverse proxies to our web services
	proxy_url, _ := url.Parse(server_configs["converterURL"])
	mux.Handle("/convert/",
		http.StripPrefix("/convert/",
			httputil.NewSingleHostReverseProxy(proxy_url)))

	proxy_url, _ = url.Parse(server_configs["aggregatorURL"])
	mux.Handle("/validate/",
		http.StripPrefix("/validate/",
			httputil.NewSingleHostReverseProxy(proxy_url)))

	// mux.HandleFunc("/api", customAPI)

	http_serv := &http.Server{
		Addr:        server_configs["http_port"],
		Handler:     mux,
		ReadTimeout: 90 * time.Second, // helps kill ghost Goroutines:
		// http://stackoverflow.com/questions/10971800/golang-http-server-leaving-open-goroutines
		//ErrorLog:   nil, // suppresses errors from stderr
	}

	log.Println("NAPLAN Validation web ui server up.")
	log.Println("Server up; listening on " + server_configs["http_port"] + "/convert")
	log.Println("Server up; listening on " + server_configs["http_port"] + "/validate")
	log.Println("Server up; listening on " + server_configs["http_port"] + "/store")

	log.Fatal(http_serv.ListenAndServe())

}
