from html.parser import HTMLParser
from Course import Schedule
import json

currentTag = ""
validTags = ["table","td","tr","th"]
courseName = ""
courseCode = ""
courseAU = ""
currentData = []
data = {}
class MyHTMLParser(HTMLParser):
    
    def handle_starttag(self, tag, attrs):
        global currentTag
        
        if(tag in validTags):
            currentTag = tag

    def handle_endtag(self, tag):
        global currentTag
        global currentData
        global data
        global courseCode, courseName, courseAU
        
        currentTag = ""
        # if len is 3: new course
        # if len is 6: new index
        # if len is 5: new entry, same index
        if(tag == "tr"):
            if(len(currentData) == 3):
                courseCode = currentData[0]
                courseName = currentData[1]
                courseAU = currentData[2]
                data[courseCode] = []
            else:
                if(len(currentData) == 5):
                    currentData = [""] + currentData
                    currentData.append("")
                elif(len(currentData) == 6):
                    if(currentData[0].isdigit()):
                        currentData.append("")
                    else:
                        currentData = [""] + currentData
                if(len(currentData) == 7):
                    data[courseCode].append(currentData)
            currentData = []

    def handle_data(self, data):
        global currentData
        
        if(currentTag != "" and data.strip() != ""):
            if(currentTag == "td"):
                currentData.append(data)

parser = MyHTMLParser()

with open("Class Schedule.html", "r", encoding='utf-8') as f:
    lines = f.readlines()
    print(len(lines))

parser.feed("".join(lines[:]))

s = []
for key, value in data.items():
    currentIndex = ""
    for val in value:
        sch = Schedule()
        if(val[0] != ""):
            currentIndex = val[0]
        sch.code = key
        sch.index = currentIndex
        sch.type = val[1]
        sch.group = val[2]
        sch.day = val[3]
        sch.time = val[4]
        sch.venue = val[5]
        sch.remark = val[6]
        s.append(sch.__dict__)

schList = json.dumps(s)

with open("schedules.json", "w") as f:
    f.write(schList)
    
#print(data["CZ1007"])
