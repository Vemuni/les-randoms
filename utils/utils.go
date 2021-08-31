package utils

import (
	"fmt"
)

var DBDateTimeFormat string
var DateFormat string
var DateTimeFormat string

func init() {
	DBDateTimeFormat = "2006-01-02 15:00:00"
	DateFormat = "02/01/2006"
	DateTimeFormat = "02/01/2006 15:04:05"
}

func HandlePanicError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
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
	return "\"" + s + "\""
}
