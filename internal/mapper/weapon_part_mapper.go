package mapper

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/iambaangkok/Card-Maker/internal/entity"
	"github.com/iambaangkok/Card-Maker/internal/reader"
)

type WeaponPartMapper interface {
	Map(csvFile reader.CSVFile) []entity.WeaponPart
}

type WeaponPartMapperImpl struct {
	ExistingEffects map[string]entity.Effect
}

func (w WeaponPartMapperImpl) Map(csvFile reader.CSVFile) []entity.WeaponPart {
	
	var weaponParts []entity.WeaponPart
	
	for row, line := range csvFile.Records {
		log.Println("mapping line", row, line)

		expectedFieldCount := reflect.TypeOf(entity.WeaponPart{}).NumField()
		if len(line) != expectedFieldCount {
			log.Fatal("invalid field count")
			continue
		}

		damage, err := strconv.Atoi(line[3])
		if err != nil { log.Fatal("damage must be int") }
		fireRate, err := strconv.Atoi(line[4])
		if err != nil { log.Fatal("fire rate must be int") }
		accuracy, err := strconv.Atoi(line[5])
		if err != nil { log.Fatal("accuracy must be int") }
		minRange, err := strconv.Atoi(line[6])
		if err != nil { log.Fatal("min range must be int") }
		maxRange, err := strconv.Atoi(line[7])
		if err != nil { log.Fatal("max range must be int") }
		ammoPerMag, err := strconv.Atoi(line[8])
		if err != nil { log.Fatal("ammo per mag must be int") }
		price, err := strconv.Atoi(line[9])
		if err != nil { log.Fatal("price must be int") }

		compatibleStrs := strings.Split(line[10], "/")
		var compatibles []entity.WeaponFrameType
		for _, compatibleStr := range compatibleStrs {
			com, exists := entity.WeaponFrameTypeNameMap[compatibleStr]
			if !exists { log.Fatal("invalid weapon frame type") }
			compatibles = append(compatibles, com)
		}

		tagStrs := strings.Split(line[11], "/")
		var tags []entity.Tag
		for _, tagStr := range tagStrs {
			tag, exists := entity.TagNameMap[tagStr]
			if !exists { log.Fatal("invalid tag") }
			tags = append(tags, tag)
		}

		effectStrs := strings.Split(line[12], "/")
		var effects []entity.Effect
		for _, effectStr := range effectStrs {
			val := strings.Split(effectStr, ":")
			effectName := val[0]
			if effectName == "-" {
				continue
			}

			effect, exists := w.ExistingEffects[effectName]
			if !exists { log.Fatal("invalid effect name") }

			if effect.HasLevel {
				if len(val) == 1 { 
					effect.Level = 1
				}else{
					level, err := strconv.Atoi(val[1])
					if err != nil {
						log.Fatal("level must be int")
					}
					effect.Level = level
				}
			}
			effects = append(effects, effect)
		}

		weaponParts = append(weaponParts, 
		entity.WeaponPart{
			Name: line[0],
			Manufacturer: line[1],
			Type: line[2],
			Damage: damage,
			FireRate: fireRate,
			Accuracy: accuracy,
			MinRange: minRange,
			MaxRange: maxRange,
			AmmoPerMag: ammoPerMag,
			Price: price,
			Compatibles: compatibles,
			Tags: tags,
			Effects: effects,
		})
	}
	return weaponParts
}