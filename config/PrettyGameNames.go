package config

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var PrettyGameNames = make(map[string]string)

// GetGameNames fetches the game display names from the gamedig github repo
func GetGameNames() {
	resp, err := http.Get("https://raw.githubusercontent.com/gamedig/node-gamedig/master/games.txt")
	if err != nil {
		log.Println(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	gameNames, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		log.Println(er2)
	}

	for _, line := range strings.Split(string(gameNames), "\n") {
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "|")
			PrettyGameNames[parts[0]] = parts[1]
		}
	}

}
