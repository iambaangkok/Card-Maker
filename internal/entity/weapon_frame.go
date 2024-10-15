package entity

import (
	"fmt"
	"reflect"
)

type WeaponFrame struct {
	Name       string
	Type       WeaponFrameType
	Damage     int
	FireRate   int
	Accuracy   int
	MaxRange   int
	AmmoPerMag int
	Price      int
	Tags       []Tag
	Effects    []Effect
}

func (w WeaponFrame) Print() {
	v := reflect.ValueOf(w)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	fmt.Println(values)
}

func (w WeaponFrame) Image() string {
	return fmt.Sprintf(`http://localhost:8081/static/img/weaponframes/%s.png`, w.Name)
}

func (w WeaponFrame) HasDamage() bool {
	return w.Damage != 0
}

func (w WeaponFrame) HasFireRate() bool {
	return w.FireRate != 0
}

func (w WeaponFrame) HasAccuracy() bool {
	return w.Accuracy != 0
}

func (w WeaponFrame) HasRange() bool {
	return w.MaxRange != 0
}

func (w WeaponFrame) HasAmmoPerMag() bool {
	return w.AmmoPerMag != 0
}

func (w WeaponFrame) GetDamageStr() string {
	sign := ""
	return fmt.Sprintf("%s%d", sign, w.Damage)
}

func (w WeaponFrame) GetFireRateStr() string {
	sign := ""
	return fmt.Sprintf("%s%d", sign, w.FireRate)
}

func (w WeaponFrame) GetAccuracyStr() string {
	sign := ""
	return fmt.Sprintf("%s%d+", sign, w.Accuracy)
}

func (w WeaponFrame) GetRangeStr() string {
	var signMin, signMax string
	signMin = ""
	signMax = ""

	return fmt.Sprintf("[%s%d, %s%d]", signMin, 0, signMax, w.MaxRange)
}

func (w WeaponFrame) GetAmmoPerMagStr() string {
	sign := ""
	return fmt.Sprintf("%s%d", sign, w.AmmoPerMag)
}