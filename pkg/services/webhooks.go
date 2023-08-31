package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/spf13/viper"
)

type DiscordWebhook struct {
	Embeds []models.Embed `json:"embeds"`
}

func SendDiscordWebhook(title, description string, fields []models.EmbedField) error {
	webhookURL := viper.GetString("webhooks.discord")
	if webhookURL == "" {
		return fmt.Errorf("Webhook URL não configurada")
	}

	embed := models.Embed{
		Title:       title,
		Description: description,
		Color:       16711680,
		Footer: &models.Footer{
			Text: "AuthAPI - Notifier",
		},
		Image: &models.Image{
			URL: "https://github.com/kevinfinalboss/StoreOps/blob/master/screenshots/Logo.jpg?raw=true",
		},
		Thumbnail: &models.Thumbnail{
			URL: "https://github.com/kevinfinalboss/StoreOps/blob/master/screenshots/Logo.jpg?raw=true",
		},
		Fields: fields,
	}

	jsonPayload, err := json.Marshal(DiscordWebhook{Embeds: []models.Embed{embed}})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Resposta de erro do Discord: %s", resp.Status)
	}

	return nil
}
