package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	gmapsKey string
	filePath string
	addr     = flag.String("addr", ":8080", "Specify the webserver port (default to 8080)")
)

func init() {
	log.SetFlags(0)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func serveHome(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("index.html")
	check(err)

	var url string
	if strings.HasPrefix(filePath, "https://") || strings.HasPrefix(filePath, "http://") {
		url = filePath
	} else {
		url = fmt.Sprintf("http://%s/%s ", req.Host, "/geo.json")
	}

	data := struct {
		GeoJSONURL string
		Key        string
		Zoom       int
	}{
		GeoJSONURL: url,
		Key:        gmapsKey,
		Zoom:       7,
	}
	err = t.Execute(w, &data)
	check(err)
}

func serveGeoJSON(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := ioutil.ReadFile(filePath)
	check(err)
	http.ServeContent(w, req, "geo.json", time.Now(), bytes.NewReader(data))
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalln("Usage: geojv /path/to/geo.json")
	}
	filePath = flag.Args()[0]
	gmapsKey = os.Getenv("GOOGLE_API_KEY")
	if gmapsKey == "" {
		log.Fatalln("GOOGLE_API_KEY environment variable not set. Please set an API key.")
	}
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/geo.json", serveGeoJSON)
	log.Printf("Serving on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
