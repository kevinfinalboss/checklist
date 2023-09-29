package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/spf13/viper"
)

type DiscordWebhook struct {
	Embeds []models.Embed `json:"embeds"`
}

func SendDiscordWebhook(title, description string, fields []models.EmbedField) error {
	return sendDiscordWebhook(title, description, fields, "")
}

func SendDiscordWebhookWithFile(title, description string, fields []models.EmbedField, filePath string) error {
	return sendDiscordWebhook(title, description, fields, filePath)
}

func sendDiscordWebhook(title, description string, fields []models.EmbedField, filePath string) error {
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

	var req *http.Request

	if filePath != "" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filePath)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}

		err = writer.WriteField("payload_json", string(jsonPayload))
		if err != nil {
			return err
		}

		err = writer.Close()
		if err != nil {
			return err
		}

		req, err = http.NewRequest("POST", webhookURL, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return err
	}

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
