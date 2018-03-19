package config

import (
	"strings"
)

const (
	SPELLCHECKER_URL = "https://www.pieceofcode.org:8080/spellcheck"
)

type ServerConfig struct {
	IsProduction    bool   `json:"is_production"`
	Port            int    `json:"port"`
	CertFile        string `json:"cert_file"`
	KeyFile         string `json:"key_file"`
	UseSpellchecker bool   `json:"use_spellchecker"`
}

type GoogleSearchConfig struct {
	SearchEngineID string `json:"search_engine_id"`
	ApiKey         string `json:"api_key"`
}

type AgentConfig struct {
	Name       string `json:"name"`
	Token      string `json:"token"`
	QueryLang  string `json:"query_lang"`
	SpeechLang string `json:"speech_lang"`
}

type DialogflowConfig struct {
	Agents []AgentConfig `json:"agents"`
}

type ExternalAgentsConfig struct {
	ClassifierUrl string `json:"classifier_url"`
}

type TestConfig struct {
	IsExactMatching	bool	`json:"is_exact_matching"`
}

func GetAgentConfigByName(agents *DialogflowConfig, agentName string) *AgentConfig {
	for index, _ := range agents.Agents {
		if strings.Compare(agents.Agents[index].Name, agentName) == 0 {
			return &agents.Agents[index]
		}
	}
	return nil
}
