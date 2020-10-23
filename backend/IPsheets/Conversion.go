package IPSheets

import (
	"strconv"
	"strings"
	"time"
)

const (
	layout  = "2006/01/02"
	layout2 = "2006-01-02"
	layout3 = "2006\\01\\02"
)

func ConvString(val interface{}) string {
	return strings.ToUpper(strings.TrimSpace(val.(string)))
}

func ConvBool(val interface{}) (bool, error) {
	return strconv.ParseBool(strings.ToUpper(strings.TrimSpace(val.(string))))
}

func ConvNum(val interface{}) (float64, error) {
	return strconv.ParseFloat(strings.ToUpper(strings.TrimSpace(val.(string))), 64)
}

func ConvDate(val interface{}) (time.Time, error) {
	val = strings.ToUpper(strings.TrimSpace(val.(string)))
	t, e := time.Parse(layout, val.(string))
	if e != nil {
		if strings.Contains(e.Error(), "-") {
			t, e = time.Parse(layout2, val.(string))
		}
		if strings.Contains(e.Error(), "\\") {
			t, e = time.Parse(layout3, val.(string))
		}
	}
	return t, e
}
