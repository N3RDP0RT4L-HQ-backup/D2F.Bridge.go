package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Config struct {
	DiscordWebhookURL string `json:"discordWebhookURL"`
	FosscordURL       string `json:"fosscordURL"`
	FosscordToken     string `json:"fosscordToken"`
}

type DiscordMessage struct {
	Content string `json:"content"`
}

type FosscordMessage struct {
	Content string `json:"content"`
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/discord", handleDiscord(config))
	http.HandleFunc("/fosscord", handleFosscord(config))
	http.ListenAndServe(":8080", nil)
}

func handleDiscord(config *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var msg DiscordMessage
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// relay message to Fosscord
		fosscordMsg := FosscordMessage{
			Content: msg.Content,
		}
		jsonValue, _ := json.Marshal(fosscordMsg)

		req, err := http.NewRequest(http.MethodPost, config.FosscordURL, bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", "Bearer "+config.FosscordToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
	}
}

func handleFosscord(config *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var msg FosscordMessage
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// relay message to Discord
		discordMsg := DiscordMessage{
			Content: msg.Content,
		}
		jsonValue, _ := json.Marshal(discordMsg)

		req, err := http.NewRequest(http.MethodPost, config.DiscordWebhookURL, bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
	}
}

func loadConfig(cfgFile string) (*Config, error) {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
