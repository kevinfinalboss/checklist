package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type SmsRequest struct {
	Key    string `json:"key"`
	Type   int    `json:"type"`
	Number string `json:"number"`
	Msg    string `json:"msg"`
}

type SmsResponse struct {
	Situacao  string `json:"situacao"`
	Codigo    string `json:"codigo"`
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}

func SendSms(number, message string) error {
	smsKey := os.Getenv("SMS_API_KEY")
	if smsKey == "" {
		return errors.New("SMS API key is missing")
	}

	smsData := SmsRequest{
		Key:    smsKey,
		Type:   9,
		Number: number,
		Msg:    message,
	}

	jsonData, err := json.Marshal([]SmsRequest{smsData})
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api.smsdev.com.br/v1/send", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var smsResp []SmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return err
	}

	if smsResp[0].Situacao != "OK" {
		return errors.New("Failed to send SMS: " + smsResp[0].Descricao)
	}

	return nil
}
