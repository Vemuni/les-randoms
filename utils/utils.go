package utils

import (
	"errors"
	"os"
)

var WEBSITE_URL string
var DBDateTimeFormat string
var DateFormat string
var DateTimeFormat string
var DebugMode bool
var TestMode bool

var NotImplementedYet error

func init() {
	WEBSITE_URL = os.Getenv("WEBSITE_URL")
	DBDateTimeFormat = "2006-01-02T15:04:05Z"
	DateFormat = "02/01/2006"
	DateTimeFormat = "02/01/2006 15:04:05"
	DebugMode = os.Getenv("DEBUG_MODE") == "TRUE"
	TestMode = os.Getenv("TEST_MODE") == "TRUE"

	NotImplementedYet = errors.New("not implemented yet")
}

func HandlePanicError(err error) {
	if err != nil {
		LogError(err.Error())
		panic(err)
	}
}

func LogNotNilError(err error) {
	if err != nil {
		LogError(err.Error())
	}
}

func UnsliceStrings(strings []string, separator string) string {
	result := ""
	for _, s := range strings {
		result = result + s + separator
	}
	return result[:len(result)-len(separator)]
}

/*
String separator : $
Escape character : !
Examples : Date$Name$Content -> {"Date", "Name", "Content"}
		   Date$ -> {"Date", ""}
		   Date!$ -> {"Date$"}
		   Date!! -> {"Date!"}

*/
func ParseDatabaseStringList(dbText string) []string {
	result := make([]string, 1)
	stringIndex := 0
	runes := []rune(dbText)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case rune('$'):
			result = append(result, "")
			stringIndex++
		case rune('!'):
			i++
			result[stringIndex] += string(dbText[i])
		default:
			result[stringIndex] += string(r)
		}
	}
	return result
}

func Esc(s string) string {
	return "'" + s + "'"
}

func FindAndRemove(list *[]string, target string) {
	for i, s := range *list {
		if s == target {
			*list = (*list)[:i+copy((*list)[i:], (*list)[i+1:])]
		}
	}
}
