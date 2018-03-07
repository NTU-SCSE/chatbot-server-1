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

func Contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item] 
    return ok
}