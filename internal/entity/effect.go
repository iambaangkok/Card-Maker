package entity

import (
	"fmt"
	"reflect"
	"strings"
)

type Effect struct {
	Name        string
	Type        string
	HasLevel    bool
	Level       int
	Description string
}

func (w Effect) Print() {
	v := reflect.ValueOf(w)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	fmt.Println(values)
}

func (e Effect) GetTypeLowercase() string {
	return strings.ToLower(e.Type)
}