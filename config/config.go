package config
import (
	"strings"
)

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

type AgentConfig struct {
	Name		string	`json:"name"`
	Token		string	`json:"token"`
	QueryLang	string	`json:"query_lang"`
	SpeechLang	string	`json:"speech_lang"`
}

type DialogflowConfig struct {
	Agents		[]AgentConfig	`json:"agents"`
}

type ExternalAgentsConfig struct {
	ClassifierUrl	string	`json:"classifier_url"`
}

func GetAgentConfigByName(agents *DialogflowConfig, agentName string) (*AgentConfig) {
	for index, _ := range(agents.Agents) {
		if(strings.Compare(agents.Agents[index].Name, agentName) == 0) {
			return &agents.Agents[index]
		}
	}
	return nil
}