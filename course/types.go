package course

type module struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	AU          int    `json:"AU"`
	PreReq      string `json:"preReq"`
	Description string `json:"description"`
}

type class struct {
	Code   string `json:"code"`
	Index  string `json:"index"`
	Type   string `json:"type"`
	Group  string `json:"group"`
	Day    string `json:"day"`
	Time   string `json:"time"`
	Venue  string `json:"venue"`
	Remark string `json:"remark"`
}
