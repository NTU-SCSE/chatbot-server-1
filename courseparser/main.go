package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Module struct {
	Name        string `json:"name"`
	PreReq      string `json:"preReq"`
	AU          string `json:"AU"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

type param struct {
	RCourseYear string `json:"r_course_yr"`
	Acad        string `json:"acad"`
	Semester    string `json:"semester"`
}

// Reference: https://stackoverflow.com/questions/24656624/golang-display-character-not-ascii-like-not-0026
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func findAndClean(re *regexp.Regexp, str string, prefixLen int, suffixLen int) string {
	result := re.FindString(str)
	resultLen := len(result)
	if resultLen == 0 {
		return ""
	}
	result = result[prefixLen : resultLen-suffixLen]
	return strings.TrimSpace(result)
}

func main() {
	fmt.Println("Start Parsing...")
	list := make([]param, 0)
	file, err := ioutil.ReadFile("./list.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(file, &list)
	modules := make([]Module, 0)

	for _, temp := range list {
		form := url.Values{
			"acadsem":     {temp.Acad + "_" + temp.Semester},
			"r_course_yr": {temp.RCourseYear},
			"acad":        {temp.Acad},
			"semester":    {temp.Semester},
			"boption":     {"CLoad"},
		}

		resp, err := http.PostForm("https://wish.wis.ntu.edu.sg/webexe/owa/AUS_SUBJ_CONT.main_display1",
			form)
		if err != nil {
			// handle error

		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		bodystr := string(body)
		newStr := strings.Replace(bodystr, "\n", "", -1)

		re := regexp.MustCompile(`<TABLE >(.*?)</TABLE>`) // get all modules
		rowRe := regexp.MustCompile(`<TR>(.*?)</TR>`)     // get headers (code, title, AU)
		columnRe := regexp.MustCompile(`<TD(.*?)>(.*?)</TD>`)
		textRe := regexp.MustCompile(`>[^>]*?</FONT>`) // get code [^>] matches any except '>'
		contentRe := regexp.MustCompile(`>[^>]*?</TD>`)

		fontSuffixLen := len("</FONT>")
		tdSuffixLen := len("</TD>")
		prereqLabel := "Prerequisite:"

		for _, match := range re.FindAllString(newStr, -1) {
			rows := rowRe.FindAllString(match, -1)
			headers := textRe.FindAllString(rows[0], -1)
			code := strings.TrimSpace(headers[0][1 : len(headers[0])-fontSuffixLen])
			title := strings.TrimSpace(headers[1][1 : len(headers[1])-fontSuffixLen])
			au := strings.TrimSpace(headers[2][1 : len(headers[2])-fontSuffixLen])

			content := findAndClean(contentRe, rows[len(rows)-1], 1, tdSuffixLen)
			prereqFound := false
			prereqs := ""
			for j := 1; j < len(rows)-1; j++ {
				columns := columnRe.FindAllString(rows[j], -1)
				label := findAndClean(textRe, columns[0], 1, fontSuffixLen)
				if len(label) > 0 && label != prereqLabel {
					prereqFound = false
				} else if prereqFound || label == prereqLabel {
					prereqFound = true
					prereq := findAndClean(textRe, columns[1], 1, fontSuffixLen)
					prereqs += " " + prereq
				}
			}

			module := Module{Name: title, Description: content, Code: code, AU: au, PreReq: prereqs}
			modules = append(modules, module)

		}
		fmt.Println(temp.RCourseYear, len(modules))
	}

	modulesJson, _ := JSONMarshal(modules, true)

	ioutil.WriteFile("courses.json", modulesJson, 0644)
}
