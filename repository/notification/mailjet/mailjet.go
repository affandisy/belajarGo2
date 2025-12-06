package mailjet

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/pobyzaarif/goshortcute"
)

type MailjetConfig struct {
	MailjetBaseURL           string
	MailjetBasicAuthUsername string
	MailjetBasicAuthPassword string
	MailjetSenderEmail       string
	MailjetSenderName        string
}

type MailjetRepository struct {
	logger        *slog.Logger
	mailjetConfig MailjetConfig
}

func NewMailjetRepository(logger *slog.Logger, cfg MailjetConfig) *MailjetRepository {
	return &MailjetRepository{
		logger,
		cfg,
	}
}

type payloadSendEmail struct {
	Messages []Messages `json:"Messages"`
}

type From struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type To struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type Messages struct {
	From     From   `json:"From"`
	To       []To   `json:"To"`
	Subject  string `json:"Subject"`
	TextPart string `json:"TextPart"`
	HTMLPart string `json:"HTMLPart"`
}

func (r *MailjetRepository) SendEmail(toName, toEmail, subject, message string) (err error) {
	url := r.mailjetConfig.MailjetBaseURL + "/v3.1/send"
	method := http.MethodPost

	toBody := []To{}
	toBody = append(toBody, To{
		Email: toEmail,
		Name:  toName,
	})

	messageBody := Messages{
		To: toBody,
		From: From{
			Email: r.mailjetConfig.MailjetSenderEmail,
			Name:  r.mailjetConfig.MailjetSenderName,
		},
		Subject:  subject,
		TextPart: message,
		HTMLPart: message,
	}
	constructMessages := []Messages{}
	constructMessages = append(constructMessages, messageBody)

	payload := payloadSendEmail{
		Messages: constructMessages,
	}

	payloadByte, _ := json.Marshal(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payloadByte)))
	if err != nil {
		return
	}

	buildBasicAuth := goshortcute.StringtoBase64Encode(r.mailjetConfig.MailjetBasicAuthUsername + ":" + r.mailjetConfig.MailjetBasicAuthPassword)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+buildBasicAuth)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}

	return fmt.Errorf("mailer service return negative response %v", res.StatusCode)
}
