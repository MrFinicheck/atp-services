package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"atp-services/internal/api"
	"atp-services/internal/app"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	dataDir := flag.String("data", "", "LevelDB data directory")
	static := flag.String("static", "frontend/dist", "Static files directory")
	flag.Parse()

	core, err := app.New(*dataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer core.Close()

	staticPath, _ := filepath.Abs(*static)
	if _, err := os.Stat(staticPath); err != nil {
		log.Printf("warning: static dir %s not found, API only mode", staticPath)
		staticPath = ""
	}

	srv := api.NewServer(core, staticPath)
	log.Printf("ATP web server on http://localhost%s (data: %s)", *addr, core.DataDir())
	if err := http.ListenAndServe(*addr, srv.Handler()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
