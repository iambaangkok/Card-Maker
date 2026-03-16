package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/iambaangkok/Card-Maker/internal/config"
	"github.com/iambaangkok/Card-Maker/internal/generic"
	"github.com/iambaangkok/Card-Maker/internal/renderer"
)

func main() {
	projectFlag := flag.String("project", "", "path to project config file (YAML/JSON)")
	flag.Parse()

	if projectFlag == nil || *projectFlag == "" {
		log.Fatal("usage: card-maker --project <path-to-project-config.yaml>")
	}

	cfg := config.LoadConfig()

	log.Printf("project config: %s", *projectFlag)
	project, err := generic.LoadProjectConfig(*projectFlag)
	if err != nil {
		log.Fatal(err)
	}

	staticImgDir := filepath.Join(project.ImageDir)
	go func() {
		http.Handle("/static/img/", http.StripPrefix("/static/img/",
			http.FileServer(http.Dir(staticImgDir))))
		if err := http.ListenAndServe("localhost:8081", nil); err != nil {
			log.Printf("static file server error: %v", err)
		}
	}()

	log.Print("waiting 2 seconds for static file server to start")
	time.Sleep(2 * time.Second)

	reg := generic.NewInMemoryTypeRegistry(project.CardTypes)
	r := renderer.ChromeRendererImpl{Config: cfg}

	if err := generic.RenderProject(project, reg, r, cfg.HTML.OutputHTMLEnabled); err != nil {
		log.Fatal(err)
	}
}
