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
	weaponFrameMapper := mapper.WeaponFrameMapperImpl{}

	csvFile := csvReader.Read("WeaponFrame")
	weaponFrames := weaponFrameMapper.Map(csvFile)
	for _, weaponFrame := range weaponFrames {
		weaponFrame.Print()
	}
}