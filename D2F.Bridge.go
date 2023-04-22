package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	DiscordWebhookURL = "https://discord.com/api/webhooks/1234567890/abcdefg" // replace with your actual Discord webhook URL
	FosscordURL       = "http://example.com/api/messages"                     // replace with your actual Fosscord API URL
	FosscordToken     = "your_api_key"                                        // replace with your actual Fosscord API key
)

type DiscordMessage struct {
	Content string `json:"content"`
}

type FosscordMessage struct {
	Content string `json:"content"`
}

func main() {
	http.HandleFunc("/discord", handleDiscord)
	http.HandleFunc("/fosscord", handleFosscord)
	http.ListenAndServe(":8080", nil)
}

func handleDiscord(w http.ResponseWriter, r *http.Request) {
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

	req, err := http.NewRequest(http.MethodPost, FosscordURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+FosscordToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
}

func handleFosscord(w http.ResponseWriter, r *http.Request) {
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

	req, err := http.NewRequest(http.MethodPost, DiscordWebhookURL, bytes.NewBuffer(jsonValue))
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
