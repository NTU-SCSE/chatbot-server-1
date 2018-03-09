package config
type ServerConfig struct {
	IsProduction	bool	`json:"is_production"`
	Port			int		`json:"port"`
	CertFile		string	`json:"cert_file"`
	KeyFile			string	`json:"key_file"`
}

type GoogleSearchConfig struct {
	SearchEngineID	string	`json:"search_engine_id"`
	ApiKey			string	`json:"api_key"`
}

type DialogflowConfig struct {
	Token		string	`json:"dialogflow_token"`
	QueryLang	string	`json:"dialogflow_query_lang"`
	SpeechLang	string	`json:"dialogflow_speech_lang"`
}