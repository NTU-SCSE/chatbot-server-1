# Courses Parser
Parsing NTU course content from https://wish.wis.ntu.edu.sg/webexe/owa/aus_subj_cont.main into JSON format.

## Overview
1. Modify list.json to include courses that will be parsed
2. Run main.go
```
go run main.go
```
3. JSON format will be printed to courses.json file

## Modifying list of courses
The following fields must be filled in for each course:
1. r_course_yr: containing course and study year information
2. acad: academic year
3. semester

The information can be obtained from browser Network Console:
```
1. Open https://wish.wis.ntu.edu.sg/webexe/owa/aus_subj_cont.main
2. Open Developer Tools (F12 in Chrome)
3. Go to "Network" tab
4. In the web interface, choose the academic year and course
5. Click "Load Content of Course(s)"
6. In the "Network", look for "AUS_SUBJ_CONT.main_display1" request
7. Look for "Form Data" section, all the information can be obtained from there
```
