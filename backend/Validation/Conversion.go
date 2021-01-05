package Validation

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "2006/01/02"
)

func ConvStringUpperNilable(val interface{}) (out string) {
	out = strings.ToUpper(strings.TrimSpace(val.(string)))
	return
}

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
	if val.(string) == "COUNT" || val.(string) == "\u2714" || val.(string) == "YES" {
		return true
	}
	if val.(string) == "NO COUNT" || val.(string) == "NO" || val.(string) == "N/A" || val.(string) == "" {
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
	processedString := strings.ToUpper(strings.TrimSpace(val.(string)))
	if processedString[0] == '$' {
		processedString = processedString[1:]
	}
	processedString = strings.ReplaceAll(processedString, ",", "")
	out, err := strconv.ParseFloat(processedString, 64)
	if err != nil {
		panic(err)
	}
	return
}

func ConvNumPos(val interface{}) (out float64) {
	if val == "" {
		panic(fmt.Errorf("%s", "&minor& nil number passes to pos only converter"))
	}
	processedString := strings.ToUpper(strings.TrimSpace(val.(string)))
	if processedString[0] == '$' {
		processedString = processedString[1:]
	}
	processedString = strings.ReplaceAll(processedString, ",", "")
	out, err := strconv.ParseFloat(processedString, 64)
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
		panic(fmt.Errorf("%s", "&minor& &continue& Backordered"))
	}
	if strings.Contains(val.(string), "CANCELLED") {
		panic(fmt.Errorf("%s", "&minor& CANCELLED"))
	}
	date := strings.ReplaceAll(val.(string), "-", "/")
	date = strings.ReplaceAll(date, "\\", "/")

	location := time.Now().Location()
	t, err := time.ParseInLocation(layout, date, location)
	if err != nil {
		if err != nil {
			panic(err)
		}
	}

	return
}
