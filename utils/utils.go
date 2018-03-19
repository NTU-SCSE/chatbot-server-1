package utils

import (
	"strconv"
	"github.com/sajari/fuzzy"
	"strings"
	"encoding/json"
	"../config"
)

func GetEnum(list []string) string {
	var result string
	for index, str := range list {
		result = result + strconv.Itoa(index+1) + ". " + str + "\t\n"
	}
	return result
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func SpellCheck(model *fuzzy.Model, query string) string {
	res := ""
	for _, word := range strings.Fields(query) {
		res += " " + model.SpellCheck(word)
	}
	return strings.TrimSpace(res)
}

func GetServerConfig(file []byte) config.ServerConfig {
	var conf config.ServerConfig
	json.Unmarshal(file, &conf)
	return conf
}

func GetDialogflowConfig(file []byte) config.DialogflowConfig {
	var conf config.DialogflowConfig
	json.Unmarshal(file, &conf)
	return conf
}

func GetExternalAgentConfig(file []byte) config.ExternalAgentsConfig {
	var conf config.ExternalAgentsConfig
	json.Unmarshal(file, &conf)
	return conf
}

func GetGoogleSearchConfig(file []byte) config.GoogleSearchConfig {
	var conf config.GoogleSearchConfig
	json.Unmarshal(file, &conf)
	return conf
}