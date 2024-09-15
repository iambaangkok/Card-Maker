package main

import (
	"github.com/iambaangkok/Card-Maker/internal/config"
	"github.com/iambaangkok/Card-Maker/internal/mapper"
	"github.com/iambaangkok/Card-Maker/internal/reader"
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


	weaponFrameCSV := csvReader.Read("WeaponFrame")
	weaponFrames := weaponFrameMapper.Map(weaponFrameCSV)
	for _, weaponFrame := range weaponFrames {
		weaponFrame.Print()
	}
}