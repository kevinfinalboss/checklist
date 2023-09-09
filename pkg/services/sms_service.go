package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/kevinfinalboss/checklist-apps/pkg/models"
)

func SendSms(number, message string) error {
	smsKey := os.Getenv("SMS_API_KEY")
	if smsKey == "" {
		return errors.New("SMS API key is missing")
	}

	smsData := models.SmsRequest{
		Key:    smsKey,
		Type:   9,
		Number: number,
		Msg:    message,
	}

	jsonData, err := json.Marshal([]models.SmsRequest{smsData})
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api.smsdev.com.br/v1/send", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var smsResp []models.SmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return err
	}

	if smsResp[0].Situacao != "OK" {
		return errors.New("Failed to send SMS: " + smsResp[0].Descricao)
	}

	return nil
}
