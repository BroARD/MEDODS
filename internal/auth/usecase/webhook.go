package usecase

import (
	"Medods/pkg/logging"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type WebHookPayload struct {
	UserID string
	IP     string
	Event  string
}

func sendWebHook(url string, whpayload WebHookPayload, logger logging.Logger) {
	jsonData, err := json.Marshal(whpayload)
	if err != nil {
		logger.Info("WebHook: ошибка преобразования payload в json")
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Info("Webhook: ошибка формирования request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second*5}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info("Webhook: ошибка связи с клиентом")
	}

	defer resp.Body.Close()
}
