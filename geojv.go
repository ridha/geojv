package main

import (
	"./web"
	"flag"
	"log"
	"os"
)

var (
	addr = flag.String("addr", ":8080", "Specify the webserver port (default to 8080)")
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalln("Usage: geojv /path/to/geo.json")
	}
	filePath := flag.Args()[0]
	gmapsKey := os.Getenv("GOOGLE_API_KEY")
	if gmapsKey == "" {
		log.Fatalln("GOOGLE_API_KEY environment variable not set. Please set an API key.")
	}
	web.StartServer(filePath, gmapsKey, *addr)
}
