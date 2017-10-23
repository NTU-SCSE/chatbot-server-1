import json
from Course import Course
codes = ["CE", "CZ", "MH"]
with open("ce.txt", "r", encoding='utf-8') as f:
    lines = f.readlines()
    print(len(lines))

courses = []
descriptionFlag = False
for i in range(len(lines)):
    # hack here
    if(lines[i][0] == '*'):
        lines[i] = lines[i][1:]
    code = lines[i][:2]
    AUIndex = lines[i].find("Acad Unit: ")
    preReqIndex = lines[i].find("Pre-requisite: ")
    if(code in codes):
        courses.append(Course())
        separator = lines[i].find(" ")
        courses[-1].code = lines[i][:separator]
        courses[-1].name = lines[i][separator+1:-1] # -1 to strip the newline
        descriptionFlag = False
    elif(AUIndex != -1):
        # hack here, CE1004 AUs is "2 / 3
        hackIndex = lines[i].find("/")
        if(hackIndex == -1):
            hackIndex = len(lines[i])
        courses[-1].AU = int(lines[i][AUIndex+len("Acad Unit: "):hackIndex])
    elif(preReqIndex != -1):
        courses[-1].preReq = lines[i][preReqIndex+len("Pre-requisite: "):]
        descriptionFlag = True
    elif(descriptionFlag):
        courses[-1].description += lines[i]

coursesList = []
for course in courses:
    coursesList.append(course.__dict__)
coursesJson = json.dumps(coursesList)


with open("ce.json", "w") as f:
    f.write(coursesJson)
        
##with open("json.txt", "r") as f:
##    s = json.load(f)
##print(s)
##print(s[0]["name"])

print(len(courses))
print(courses[0].code)
print(courses[0].name)
print(courses[0].AU)
print(courses[0].preReq)
print(courses[0].description)
