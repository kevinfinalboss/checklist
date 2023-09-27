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
		err := fmt.Errorf("Webhook URL não configurada")
		fmt.Println("Erro no SendDiscordWebhook:", err)
		return err
	}

	embed := models.Embed{
		Title:       title,
		Description: description,
		Fields:      fields,
	}

	jsonPayload, err := json.Marshal(DiscordWebhook{Embeds: []models.Embed{embed}})
	if err != nil {
		fmt.Println("Erro ao serializar o payload do webhook:", err)
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Erro ao criar nova requisição HTTP:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição do webhook:", err)
		return err
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Corpo da Resposta:", string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		err := fmt.Errorf("Resposta de erro do Discord: %s", resp.Status)
		fmt.Println("Erro no SendDiscordWebhook:", err)
		return err
	}

	return nil
}
