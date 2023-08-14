package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
)

type DiscordWebhook struct {
	Embeds []models.Embed `json:"embeds"`
}

func SendDiscordWebhook(title, description string) error {
	webhookURL := "https://discord.com/api/webhooks/1103840164062707762/hmu05z5RrS4ya4QTHBKT7XxSaCfS1JxoACWZ750lzje0sZpejBY_6tu0AzK1pAshzJ4m"

	embed := models.Embed{
		Title:       title,
		Description: description,
		Color:       16711680,
		Footer: &models.Footer{
			Text: "CheckList API - Error Notifier",
		},
	}

	payload := DiscordWebhook{
		Embeds: []models.Embed{embed},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		utils.Logger.Error("Erro no payload:", err)
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		utils.Logger.Error("Error ao enviar mensagem ao webhook:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		utils.Logger.Error("Resposta de erro do Discord:", resp.Status)
		return fmt.Errorf("Resposta de erro do Discord: %s", resp.Status)
	}

	fmt.Println("Discord webhook enviada com sucesso")
	return nil
}
