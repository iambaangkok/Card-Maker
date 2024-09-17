package entity

import (
	"fmt"
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
}

func (w WeaponPart) Print() {
	v := reflect.ValueOf(w)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	fmt.Println(values)
}
