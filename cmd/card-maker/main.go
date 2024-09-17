package main

import (
	"log"

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


	weaponFrameCSV := csvReader.Read("WeaponFrame")
	weaponFrames := weaponFrameMapper.Map(weaponFrameCSV)
	for _, weaponFrame := range weaponFrames {
		weaponFrame.Print()
	}

	weaponPartCSV := csvReader.Read("WeaponPart")
	weaponParts := weaponPartMapper.Map(weaponPartCSV)
	for _, weaponPart := range weaponParts {
		weaponPart.Print()
	}

	// test chromedp
	err := renderer.RenderElement("div.title-container", "test-render.png")
	if err != nil {
		log.Fatal(err)
	}
}

