package mapper

import (
	"log"
	"reflect"
	"strconv"

	"github.com/iambaangkok/Card-Maker/internal/entity"
	"github.com/iambaangkok/Card-Maker/internal/reader"
)

type WeaponFrameMapper interface {
	Map(csvFile reader.CSVFile) []entity.WeaponFrame
}

type WeaponFrameMapperImpl struct {
}

func (w WeaponFrameMapperImpl) Map(csvFile reader.CSVFile) []entity.WeaponFrame {
	
	var weaponFrames []entity.WeaponFrame
	
	for row, line := range csvFile.Records {
		log.Println("mapping line", row, line)

		expectedFieldCount := reflect.TypeOf(entity.WeaponFrame{}).NumField()
		if len(line) <= expectedFieldCount {
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
		// tags := line[8]
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
		})
	}
	return weaponFrames
}