package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
	"strings"
)

var URL = os.Getenv("SLACK_URL")

type EmailMessage struct {
	Mail    events.SimpleEmailMessage `json:"mail"`
	Content string                    `json:"content"`
}

type SlackData struct {
	Blocks []SlackBlockData `json:"blocks"`
}

type SlackBlockData struct {
	Type string        `json:"type"`
	Text SlackTextData `json:"text"`
}

type SlackTextData struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func formatSlackBlock(from []string, subject string, content string) SlackBlockData {
	var res SlackBlockData
	res.Type = "section"
	res.Text.Type = "mrkdwn"
	res.Text.Text = fmt.Sprintf("*From:* %s\n*Subject*: %s\n%s", strings.Join(from, " "), subject, content)
	return res
}

func sendToSlack(blocks []SlackBlockData) error {
	data := SlackData{
		Blocks: blocks,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(b))

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid response")
	}

	return nil
}

func getMessageContent(encoded string) (string, error) {
	message, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	messageStr := string(message)

	idx := strings.Index(messageStr, "\r\n\r\n")
	if idx < 0 {
		return "", errors.New("Invalid message format")
	}

	messageStr = messageStr[idx:]
	messageStr = strings.Replace(messageStr, "\r\n", "\n", -1)
	return messageStr, nil
}

func HandleRequest(ctx context.Context, event events.SNSEvent) error {
	blocks := make([]SlackBlockData, 0)
	for _, record := range event.Records {
		snsMsg := record.SNS

		var email EmailMessage
		if err := json.Unmarshal([]byte(snsMsg.Message), &email); err != nil {
			return err
		}

		message, err := getMessageContent(email.Content)
		if err != nil {
			return err
		}
		block := formatSlackBlock(email.Mail.CommonHeaders.From, email.Mail.CommonHeaders.Subject, message)
		blocks = append(blocks, block)
	}

	return sendToSlack(blocks)
}

func main() {
	lambda.Start(HandleRequest)
}
