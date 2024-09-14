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
	// Effects    []Effects
}

func (w WeaponFrame) Print() {
	v := reflect.ValueOf(w)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	fmt.Println(values)
}