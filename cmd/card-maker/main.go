package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/iambaangkok/Card-Maker/internal/config"
	"github.com/iambaangkok/Card-Maker/internal/mapper"
	"github.com/iambaangkok/Card-Maker/internal/reader"
	"github.com/iambaangkok/Card-Maker/internal/renderer"
)

func main() {
	cfg := config.LoadConfig()
	
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

	// serve static files for html
	go func () {
		http.Handle("/static/img/", http.StripPrefix("/static/img/",
		http.FileServer(http.Dir(path.Join("./internal/template", "/img/")))))
		http.ListenAndServe("localhost:8081", nil)
	}()

	// load html template
	template, err := template.ParseFiles("./internal/template/html/WeaponFrame.html")
	if err != nil {
		log.Fatal("unable to read html template")
	}

	for _, weaponFrame := range weaponFrames {
		data := weaponFrame

		var buf bytes.Buffer
		err = template.Execute(&buf, data)
		if err != nil {
			log.Fatal("error while applying html")
		}
		parsedHtml := buf.String()
		outputDir := "./output/"
		outputFilePath := outputDir + "WeaponFrame_" + weaponFrame.Name + ".html"
		err := os.WriteFile(outputFilePath, []byte(parsedHtml), 0644)
		if err != nil {
			log.Fatal("error while saving parsed html")
		}
		// render to png
		outputFilePath = "WeaponFrame_" + weaponFrame.Name + ".png"
		err = renderer.RenderHTMLToPNG(parsedHtml, outputFilePath)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("rendered ", outputFilePath)
	}


	

}

