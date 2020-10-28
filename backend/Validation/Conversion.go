package Validation

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	layout  = "2006/01/02"
	layout2 = "2006-01-02"
	layout3 = "2006\\01\\02"
)

func ConvStringUpper(val interface{}) string {
	return strings.ToUpper(strings.TrimSpace(val.(string)))
}

func ConvString(val interface{}) string {
	return strings.TrimSpace(val.(string))
}

func ConvBool(val interface{}) (out bool) {
	val = strings.ToUpper(strings.TrimSpace(val.(string)))
	if val.(string) == "COUNT" || val.(string) == "\u2714" {
		return true
	}
	if val.(string) == "" {
		return false
	}
	out, err := strconv.ParseBool(val.(string))
	if err != nil {
		panic(err)
	}
	return
}

func ConvNum(val interface{}) (out float64) {
	if val == "" {
		return 0
	}
	out, err := strconv.ParseFloat(strings.ToUpper(strings.TrimSpace(val.(string))), 64)
	if err != nil {
		panic(err)
	}
	return
}

func ConvDate(val interface{}) (t time.Time) {
	val = strings.ToUpper(strings.TrimSpace(val.(string)))
	if val == "" {
		return
	}
	if strings.Contains(val.(string), "B") {
		panic(fmt.Errorf("%s", "&minor& Backordered"))
	}
	t, err := time.Parse(layout, val.(string))
	if err != nil {
		if strings.Contains(err.Error(), "-") {
			t, err = time.Parse(layout2, val.(string))
		}
		if strings.Contains(err.Error(), "\\") {
			t, err = time.Parse(layout3, val.(string))
		}
		if err != nil {
			panic(err)
		}
	}

	return
}
