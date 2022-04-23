package gamequery

import (
	"GoGSM/config"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
)

type Query struct {
	Name       string `json:"name,omitempty"`
	Map        string `json:"map,omitempty"`
	Password   bool   `json:"password,omitempty"`
	Desc       string `json:"desc,omitempty"`
	MaxPlayers int    `json:"maxplayers,omitempty"`
	Players    int    `json:"players,omitempty"`
	Bots       int    `json:"bots,omitempty"`
	Connect    string `json:"connect,omitempty"`
	Ping       int    `json:"ping,omitempty"`
	Error      string `json:"error,omitempty"`
}

type GameDigPlayers struct {
	Name string
}

var ErrServerNotFound = errors.New("server not found")

// ServerQuery returns a Query struct with the server information, or an error if the server is not found
func ServerQuery(query *config.Server) (Query, error) {
	cmd := exec.Command("gamedig", "--type", query.Game, "--host", query.Hostname, "--port", strconv.Itoa(query.Port))
	out, err := cmd.Output()
	if err != nil {
		return Query{}, err
	}
	var rawJsonData map[string]interface{}
	er2 := json.Unmarshal(out, &rawJsonData)
	if er2 != nil {
		return Query{}, er2
	}

	if rawJsonData["error"] == "Failed all 1 attempts" {
		return Query{}, ErrServerNotFound
	}

	var response Query

	if rawJsonData["name"] != nil {
		response.Name = rawJsonData["name"].(string)
	}

	if rawJsonData["map"] != nil {
		response.Map = rawJsonData["map"].(string)
	}

	if rawJsonData["password"] == true {
		response.Password = true
	} else {
		response.Password = false
	}

	if rawJsonData["desc"] != nil {
		response.Desc = rawJsonData["desc"].(string)
	}

	if rawJsonData["maxplayers"] != nil {
		response.MaxPlayers = int(rawJsonData["maxplayers"].(float64))
	}

	switch rawJsonData["players"].(type) {
	case []interface{}:
		response.Players = len(rawJsonData["players"].([]interface{}))
	case int:
		response.Players = int(rawJsonData["players"].(float64))
	}

	if rawJsonData["bots"] != nil {
		response.Bots = len(rawJsonData["bots"].([]interface{}))
	}

	if rawJsonData["connect"] != nil {
		response.Connect = rawJsonData["connect"].(string)
	}

	if rawJsonData["ping"] != nil {
		response.Ping = int(rawJsonData["ping"].(float64))
	}

	if rawJsonData["error"] != nil {
		response.Error = rawJsonData["error"].(string)
	}

	return response, nil
}
