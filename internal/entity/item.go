package entity

import (
	"fmt"
	"html/template"
	"reflect"
)

type Item struct {
	Name       string
	UsageLimit int
	Price	   int
	Effects    []Effect
	ImgTag	 	 template.HTML
}

func (w Item) Print() {
	v := reflect.ValueOf(w)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	fmt.Println(values)
}

func (w Item) Image() string {
	return fmt.Sprintf(`http://localhost:8081/static/img/items/%s.png`, w.Name)
}

func (w Item) HasUsageLimit() bool {
	return w.UsageLimit > 0
}

func (w Item) GetUsageLimitArray() []int {
	if w.UsageLimit == 0 {
		return make([]int, 0)
	}
	arr := make([]int, w.UsageLimit)
	for i := 0; i < w.UsageLimit; i++ {
		arr[i] = i+1
	}
	return arr
}