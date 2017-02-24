package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleHomePage(filePath, gmapsKey string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		page, err := Asset("index.html")
		check(err)

		t := template.Must(template.New("index.html").Parse(string(page)))

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
}

func handleGeoJSON(filePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := ioutil.ReadFile(filePath)
		check(err)
		http.ServeContent(w, req, "geo.json", time.Now(), bytes.NewReader(data))
	}
}

func StartServer(filePath string, gmapsKey string, addr string) {
	http.HandleFunc("/", handleHomePage(filePath, gmapsKey))
	http.HandleFunc("/geo.json", handleGeoJSON(filePath))
	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
