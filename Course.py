class Course(object):
    def __init__(self):
        self.code = ""
        self.name = ""
        self.AU = 0
        self.preReq = []
        self.description = ""

class Schedule(object):
    def __init__(self):
        self.code = ""
        self.index = ""
        self.group = ""
        self.type = ""
        self.day = ""
        self.time = ""
        self.venue = ""
        self.remark = ""

