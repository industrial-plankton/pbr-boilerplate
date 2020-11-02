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

func ConvStringUpper(val interface{}) (out string) {
	out = strings.ToUpper(strings.TrimSpace(val.(string)))
	if out == "" {
		panic(fmt.Errorf("%s", "&minor& Nil string passed"))
	}
	return
}

func ConvString(val interface{}) (out string) {
	out = strings.TrimSpace(val.(string))
	if out == "" {
		panic(fmt.Errorf("%s", "&minor& Nil string passed"))
	}
	return
}

func ConvBool(val interface{}) (out bool) {
	val = strings.ToUpper(strings.TrimSpace(val.(string)))
	if val.(string) == "COUNT" || val.(string) == "\u2714" {
		return true
	}
	if val.(string) == "NO COUNT" {
		return false
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

func ConvNumPos(val interface{}) (out float64) {
	if val == "" {
		panic(fmt.Errorf("%s", "&minor& nil number passes to pos only converter"))
	}
	out, err := strconv.ParseFloat(strings.ToUpper(strings.TrimSpace(val.(string))), 64)
	if err != nil {
		panic(err)
	}
	if out < 0 {
		panic(fmt.Errorf("%s", "&minor& negative number passes to pos only converter"))
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
