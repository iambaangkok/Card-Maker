package mapper

import (
	"log"
	"strconv"

	"github.com/iambaangkok/Card-Maker/internal/entity"
	"github.com/iambaangkok/Card-Maker/internal/reader"
)

type EffectMapper interface {
	Map(csvFile reader.CSVFile) []entity.Effect
}

type EffectMapperImpl struct {
}

func (w EffectMapperImpl) MapToMap(csvFile reader.CSVFile) map[string]entity.Effect{
	effectsNameMap := make(map[string]entity.Effect)
	
	effectList := w.MapToList(csvFile)
	for _, e := range effectList {
		effectsNameMap[e.Name] = e
	}

	return effectsNameMap
}

func (w EffectMapperImpl) MapToList(csvFile reader.CSVFile) []entity.Effect {
	
	var effects []entity.Effect
	
	for row, line := range csvFile.Records {
		log.Println("mapping line", row, line)

		hasLevelInt, err := strconv.Atoi(line[2])
		if err != nil { log.Fatal("has level must be int") }
		hasLevel := hasLevelInt != 0

		effects = append(effects, 
		entity.Effect{
			Name: line[0],
			Type: line[1],
			HasLevel: hasLevel,
			Level: 0,
			Description: line[3],
		})
	}
	return effects
}