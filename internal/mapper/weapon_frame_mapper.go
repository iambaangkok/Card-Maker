package mapper

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/iambaangkok/Card-Maker/internal/entity"
	"github.com/iambaangkok/Card-Maker/internal/reader"
)

type WeaponFrameMapper interface {
	Map(csvFile reader.CSVFile) []entity.WeaponFrame
}

type WeaponFrameMapperImpl struct {
	ExistingEffects map[string]entity.Effect
}

func (w WeaponFrameMapperImpl) Map(csvFile reader.CSVFile) []entity.WeaponFrame {
	
	var weaponFrames []entity.WeaponFrame
	
	for row, line := range csvFile.Records {
		log.Println("mapping line", row, line)

		expectedFieldCount := reflect.TypeOf(entity.WeaponFrame{}).NumField()
		if len(line) != expectedFieldCount {
			log.Fatal("invalid field count")
			continue
		}

		weaponFrameType, exists := entity.WeaponFrameTypeNameMap[line[1]]
		if !exists { log.Fatal("invalid weapon frame type") }
		damage, err := strconv.Atoi(line[2])
		if err != nil { log.Fatal("damage must be int") }
		fireRate, err := strconv.Atoi(line[3])
		if err != nil { log.Fatal("fire rate must be int") }
		accuracy, err := strconv.Atoi(line[4])
		if err != nil { log.Fatal("accuracy must be int") }
		maxRange, err := strconv.Atoi(line[5])
		if err != nil { log.Fatal("max range must be int") }
		ammoPerMag, err := strconv.Atoi(line[6])
		if err != nil { log.Fatal("ammo per mag must be int") }
		price, err := strconv.Atoi(line[7])
		if err != nil { log.Fatal("price must be int") }
		tagStrs := strings.Split(line[8], "/")
		var tags []entity.Tag
		for _, tagStr := range tagStrs {
			tag, exists := entity.TagNameMap[tagStr]
			if !exists { log.Fatal("invalid tag") }
			tags = append(tags, tag)
		}
		effectStrs := strings.Split(line[9], "/")
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

		weaponFrames = append(weaponFrames, 
		entity.WeaponFrame{
			Name: line[0],
			Type: weaponFrameType,
			Damage: damage,
			FireRate: fireRate,
			Accuracy: accuracy,
			MaxRange: maxRange,
			AmmoPerMag: ammoPerMag,
			Price: price,
			Tags: tags,
			Effects: effects,
		})
	}
	return weaponFrames
}