package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/iambaangkok/Card-Maker/internal/config"
	"github.com/iambaangkok/Card-Maker/internal/generic"
	"github.com/iambaangkok/Card-Maker/internal/mapper"
	"github.com/iambaangkok/Card-Maker/internal/reader"
	"github.com/iambaangkok/Card-Maker/internal/renderer"
)

func main() {
	projectFlag := flag.String("project", "", "path to project config file (YAML/JSON) for generic mode")
	flag.Parse()

	cfg := config.LoadConfig()

	// If --project is provided, use the new generic engine path.
	if projectFlag != nil && *projectFlag != "" {
		runGeneric(cfg, *projectFlag)
		return
	}

	csvReader := reader.CSVReaderImpl{
		Config: cfg,
	}
	effectMapper := mapper.EffectMapperImpl{}

	effectsCSV := csvReader.Read("Effects")
	effectsNameMap := effectMapper.MapToMap(effectsCSV)

	weaponFrameMapper := mapper.WeaponFrameMapperImpl{
		ExistingEffects: effectsNameMap,
	}
	weaponPartMapper := mapper.WeaponPartMapperImpl{
		ExistingEffects: effectsNameMap,
	}
	itemMapper := mapper.ItemMapperImpl{
		ExistingEffects: effectsNameMap,
	}
	renderer := renderer.ChromeRendererImpl{
		Config: cfg,
	}

	// map weapon frame
	weaponFrameCSV := csvReader.Read("WeaponFrame")
	weaponFrames := weaponFrameMapper.Map(weaponFrameCSV)
	for _, weaponFrame := range weaponFrames {
		weaponFrame.Print()
	}

	// map weapon part
	weaponPartCSV := csvReader.Read("WeaponPart")
	weaponParts := weaponPartMapper.Map(weaponPartCSV)
	for _, weaponPart := range weaponParts {
		weaponPart.Print()
	}
	
	// map item
	itemCSV := csvReader.Read("Item")
	items := itemMapper.Map(itemCSV)
	for _, item := range items {
		item.Print()
	}

	// serve static files for html
	go func () {
		http.Handle("/static/img/", http.StripPrefix("/static/img/",
		http.FileServer(http.Dir(path.Join("./internal/template", "/img/")))))
		http.ListenAndServe("localhost:8081", nil)
	}()

	// WEAPON FRAMES
	if cfg.Renderer.RenderWeaponFrameEnabled {
		// load html tem
		tem, err := template.ParseFiles("./internal/template/html/WeaponFrame.html")
		if err != nil {
			log.Fatal("unable to read html template")
		}
	
		for _, weaponFrame := range weaponFrames {
			data := weaponFrame
	
			var buf bytes.Buffer
			err = tem.Execute(&buf, data)
			if err != nil {
				log.Fatal("error while applying html")
			}
			parsedHtml := buf.String()
			if cfg.HTML.OutputHTMLEnabled {
				outputDir := "./output/"
				outputFilePath := outputDir + "WeaponFrame_" + weaponFrame.Name + ".html"
				err := os.WriteFile(outputFilePath, []byte(parsedHtml), 0644)
				if err != nil {
					log.Fatal("error while saving parsed html")
				}
			}
			// render to png
			outputFilePath := "WeaponFrame_" + weaponFrame.Name + ".png"
			err = renderer.RenderHTMLToPNG(parsedHtml, outputFilePath)
			if err != nil {
				log.Fatal(err)
			}
	
			log.Print("rendered ", outputFilePath)
		}
	}
	

	// WEAPON PARTS
	if cfg.Renderer.RenderWeaponPartEnabled {
		// load html template
		tem, err := template.ParseFiles("./internal/template/html/WeaponPart.html")
		if err != nil {
			log.Fatal("unable to read html template")
		}

		for _, weaponPart := range weaponParts {
			data := weaponPart

			var buf bytes.Buffer
			err = tem.Execute(&buf, data)
			if err != nil {
				log.Fatal("error while applying html")
			}
			parsedHtml := buf.String()
			if cfg.HTML.OutputHTMLEnabled {
				outputDir := "./output/"
				outputFilePath := outputDir + "WeaponPart_" + weaponPart.Name + ".html"
				err := os.WriteFile(outputFilePath, []byte(parsedHtml), 0644)
				if err != nil {
					log.Fatal("error while saving parsed html")
				}
			}
			
			// render to png
			outputFilePath := "WeaponPart_" + weaponPart.Name + ".png"
			err = renderer.RenderHTMLToPNG(parsedHtml, outputFilePath)
			if err != nil {
				log.Fatal(err)
			}

			log.Print("rendered ", outputFilePath)
		}
	}
	// ITEMS
	if cfg.Renderer.RenderItemEnabled {
		// load html tem
		tem, err := template.ParseFiles("./internal/template/html/Item.html")
		if err != nil {
			log.Fatal("unable to read html template")
		}

		for _, item := range items {
			data := item

			var buf bytes.Buffer
			err = tem.Execute(&buf, data)
			if err != nil {
				log.Fatal("error while applying html")
			}
			parsedHtml := buf.String()
			if cfg.HTML.OutputHTMLEnabled {
				outputDir := "./output/"
				outputFilePath := outputDir + "Item_" + item.Name + ".html"
				err := os.WriteFile(outputFilePath, []byte(parsedHtml), 0644)
				if err != nil {
					log.Fatal("error while saving parsed html")
				}
			}
			// render to png
			outputFilePath := "Item_" + item.Name + ".png"
			err = renderer.RenderHTMLToPNG(parsedHtml, outputFilePath)
			if err != nil {
				log.Fatal(err)
			}

			log.Print("rendered ", outputFilePath)
		}
	}
}

func runGeneric(cfg config.Config, projectConfigPath string) {
	log.Printf("running in generic mode with project config %s", projectConfigPath)

	project, err := generic.LoadProjectConfig(projectConfigPath)
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

	r := renderer.ChromeRendererImpl{
		Config: cfg,
	}

	if err := generic.RenderProject(project, reg, r, cfg.HTML.OutputHTMLEnabled); err != nil {
		log.Fatal(err)
	}
}


