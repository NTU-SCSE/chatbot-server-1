package utils
import (
	"strconv"
)
func GetEnum(list []string) string {
	var result string
	for index, str := range list {
		result = result + strconv.Itoa(index + 1) + ". " + str + "\t\n"
	}
	return result
}