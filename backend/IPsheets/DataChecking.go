package IPSheets

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func First5(f func(interface{}), data [][]interface{}) {
	size := len(data)
	if size > 5 {
		size = 5
	}

	for i := 0; i < size; i++ {
		for _, e := range data[i] {
			f(e)
		}
	}
}

func CheckTypes(data interface{}) {
	xType := reflect.TypeOf(data)
	xValue := reflect.ValueOf(data)
	fmt.Println(xType, xValue)
}

func Dates(data interface{}) {
	t, e := time.Parse(layout, data.(string))
	if e != nil {
		if strings.Contains(e.Error(), "-") {
			t, e = time.Parse(layout2, data.(string))
		}
		if strings.Contains(e.Error(), "\\") {
			t, e = time.Parse(layout3, data.(string))
		}
	}
	fmt.Println(t, data, e)
}

func TryNum(data interface{}) {
	d, e := strconv.ParseFloat(data.(string), 64)
	fmt.Println(d, data, e)
}

func TryBool(data interface{}) {
	d, e := strconv.ParseBool(data.(string))
	fmt.Println(d, data, e)
}

func Printmap(data map[string]float32) {
	for key, value := range data {
		fmt.Println(key, ",", value)
	}
}
