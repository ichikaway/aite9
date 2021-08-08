package notification

import (
	"aite9/printer"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

func createSlackMessage(title, text string) []byte {
	slackUsername, ok2 := os.LookupEnv("SLACK_USERNAME")
	if !ok2 {
		slackUsername = ""
	}
	slackChannel, ok3 := os.LookupEnv("SLACK_CHANNEL")
	if !ok3 {
		slackChannel = ""
	}
	iconEmoji, ok4 := os.LookupEnv("SLACK_ICON_EMOJI")
	if !ok4 {
		iconEmoji = ""
	}
	iconURL, ok5 := os.LookupEnv("SLACK_ICON_URL")
	if !ok5 {
		iconURL = ""
	}
	freeText, ok6 := os.LookupEnv("SLACK_FREE_TEXT")
	if !ok6 {
		freeText = ""
	}

	type attachments struct {
		Color string `json:"color"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	type slack struct {
		Username     string        `json:"username"`
		IconEmoji    string        `json:"icon_emoji"`
		IconURL      string        `json:"icon_url"`
		Channel      string        `json:"channel"`
		Text         string        `json:"text"`
		Attachements []attachments `json:"attachments"`
	}

	webhooks := slack{
		Username:  slackUsername,
		IconEmoji: iconEmoji,
		IconURL:   iconURL,
		Channel:   slackChannel,
		Text:      freeText,
		Attachements: []attachments{
			{
				Color: "warning",
				Title: title,
				Text:  text,
			},
		},
	}

	params, _ := json.Marshal(webhooks)
	return params
}

func PostSlack(title, text string) {
	slackWebhookURL, ok1 := os.LookupEnv("SLACK_WEBHOOK_URL")

	if !ok1 {
		return
	}

	printer.Printf("Sending warning message to slack.\n")

	params := createSlackMessage(title, text)

	resp, err := http.PostForm(
		slackWebhookURL,
		url.Values{"payload": {string(params)}},
	)
	if err == nil {
		defer resp.Body.Close()
	}
}
