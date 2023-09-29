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
		webhookURL = "https://discord.com/api/webhooks/1157394403594346556/MypTfWBNLwQ4Cwgn0YtPza_HwBw9WvfFCV1T6IXPJwgALwOJqqJu2Rr2Z_gvZwxN7ATg" // URL de fallback
		fmt.Println("Usando URL de Webhook de fallback.")
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

	fmt.Println("Payload JSON:", string(jsonPayload)) // Log do payload

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
	fmt.Println("Cabeçalhos da Resposta:", resp.Header)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		err := fmt.Errorf("Resposta de erro do Discord: %s", resp.Status)
		fmt.Println("Erro no SendDiscordWebhook:", err)
		return err
	}

	return nil
}
