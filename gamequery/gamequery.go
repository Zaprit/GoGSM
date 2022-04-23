package gamequery

import (
	"GoGSM/config"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
)

type Query struct {
	Name       string           `json:"name,omitempty"`
	Map        string           `json:"map,omitempty"`
	Password   bool             `json:"password,omitempty"`
	Desc       string           `json:"desc,omitempty"`
	MaxPlayers int              `json:"maxplayers,omitempty"`
	Players    []GameDigPlayers `json:"players,omitempty"`
	Bots       []GameDigPlayers `json:"bots,omitempty"`
	Connect    string           `json:"connect,omitempty"`
	Ping       int              `json:"ping,omitempty"`
	Error      string           `json:"error,omitempty"`
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
	var response Query
	er2 := json.Unmarshal(out, &response)
	if er2 != nil {
		return Query{}, er2
	}

	if response.Error == "Failed all 1 attempts" {
		return Query{}, ErrServerNotFound
	}

	return response, nil
}
