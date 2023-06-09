package utils

import (
	"fmt"
	"strconv"
	"time"
)

func StringToBool(value string, _default bool) bool {

	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		return _default
	}
}

func StringToFloat(value string) float64 {

	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return float
}

func StringToInt(value string) int64 {

	integer, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return integer
}

func MapStringsToString(queryMap map[string]interface{}) string {

	var query string = "?"
	for key, value := range queryMap {

		stringValue := fmt.Sprintf("%+v", value)
		if stringValue == "" {
			continue
		}

		query += (key + "=" + stringValue + "&")
	}
	return query[:len(query)-1]
}

func MapStringsAppend(queryIn string, queryMap map[string]interface{}) string {
	var query string = queryIn + "&"
	for key, value := range queryMap {

		stringValue := fmt.Sprintf("%+v", value)
		if stringValue == "" {
			continue
		}

		query += (key + "=" + stringValue + "&")
	}
	return query[:len(query)-1]

}

func ParseStringToTime(dateStr string) time.Time {

	date, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		fmt.Println(err)
		var t time.Time
		return t
	}
	return date
}
