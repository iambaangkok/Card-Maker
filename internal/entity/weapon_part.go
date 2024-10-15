package entity

import (
	"fmt"
	"html/template"
	"reflect"
)

type WeaponPart struct {
	Name         string
	Manufacturer string
	Type         string
	Damage       int
	FireRate     int
	Accuracy     int
	MinRange     int
	MaxRange     int
	AmmoPerMag   int
	Price        int
	Compatibles  []WeaponFrameType
	Tags         []Tag
	Effects      []Effect
	ImgTag	 	 template.HTML
}

func (w WeaponPart) Print() {
	v := reflect.ValueOf(w)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	fmt.Println(values)
}

func (w WeaponPart) Image() string {
	return fmt.Sprintf(`http://localhost:8081/static/img/weaponparts/%s.png`, w.Name)
}

func (w WeaponPart) HasDamage() bool {
	return w.Damage != 0
}

func (w WeaponPart) HasFireRate() bool {
	return w.FireRate != 0
}

func (w WeaponPart) HasAccuracy() bool {
	return w.Accuracy != 0
}

func (w WeaponPart) HasRange() bool {
	return w.MinRange != 0 || w.MaxRange != 0
}

func (w WeaponPart) HasAmmoPerMag() bool {
	return w.AmmoPerMag != 0
}

func (w WeaponPart) GetDamageStr() string {
	sign := ""
	if w.Damage > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%d", sign, w.Damage)
}

func (w WeaponPart) GetFireRateStr() string {
	sign := ""
	if w.FireRate > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%d", sign, w.FireRate)
}

func (w WeaponPart) GetAccuracyStr() string {
	sign := ""
	if w.Accuracy > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%d", sign, w.Accuracy)
}

func (w WeaponPart) GetRangeStr() string {
	var signMin, signMax string
	if w.MinRange > 0 {
		signMin = "+"
	} else if w.MinRange < 0 {
		signMin = ""
	}
	if w.MaxRange > 0 {
		signMax = "+"
	} else if w.MaxRange < 0 {
		signMax = ""
	}

	return fmt.Sprintf("[%s%d, %s%d]", signMin, w.MinRange, signMax, w.MaxRange)
}

func (w WeaponPart) GetAmmoPerMagStr() string {
	sign := ""
	if w.AmmoPerMag > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%d", sign, w.AmmoPerMag)
}

