package course

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Course struct {
	Modules []module
	Classes []class
}

func NewCourse() *Course {
	c := Course{Modules: make([]module, 0), Classes: make([]class, 0)}
	file, err := ioutil.ReadFile("./cs.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(file, &c.Modules)

	var CECourses []module
	file, err = ioutil.ReadFile("./ce.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(file, &CECourses)
	c.Modules = append(c.Modules, CECourses...)

	// Get the data of course schedules and venues
	file, err = ioutil.ReadFile("./schedules.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(file, &c.Classes)
	return &c
}

func ParseCourseCode(param string) string {
	var validCode = regexp.MustCompile(`^c[z|e](/c[z|e])?[0-9]{4}$`)
	return validCode.FindString(param)
}

func (c *Course) GetCourseCode(param string) (string, string) {
	param = strings.TrimSpace(param)

	if _, err := strconv.Atoi(param[len(param)-4:]); err == nil {
		return param, ""
	}
	result := ""
	auxResult := ""

	for _, mod := range c.Modules {
		if strings.ToLower(strings.TrimSpace(mod.Name)) == strings.ToLower(param) {
			if result == "" {
				result = mod.Code
			} else if mod.Code[2] != '/' {
				auxResult = "CE/CZ" + mod.Code[2:]
			}
		}
	}
	return result, auxResult
}

func GetSchedulePrint(param class) string {
	var result string
	if param.Type == "LEC/STUDIO" {
		result = "Lecture"
	} else if param.Type == "TUT" {
		result = "Tutorial"
	} else {
		result = param.Type
	}
	result = result + " on " + param.Day + ", " + param.Time + " at " + param.Venue
	return result
}

func (c *Course) getIndex(code string) map[string]bool {
	result := map[string]bool{}
	for _, class := range c.Classes {
		if strings.ToLower(strings.TrimSpace(class.Code)) == strings.ToLower(code) {
			result[class.Index] = true
		}
	}
	return result
}

func (c *Course) GetIndexString(code string) string {
	indexList := c.getIndex(code)
	indexStrings := []string{}
	for key, _ := range indexList {
		indexStrings = append(indexStrings, key)
	}
	sort.Strings(indexStrings)
	result := ""
	for _, key := range indexStrings {
		result = result + key + "\n"
	}
	return result
}
