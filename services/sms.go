package services

import (
	"fmt"
	"log"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendSMS sends an OTP to a given phone number via Twilio
func SendSMS(to string, code string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	if accountSid == "" || authToken == "" {
		log.Println("Twilio credentials not set, mocking SMS send")
		log.Printf("MOCK SMS to %s: Your OTP is %s\n", to, code)
		return nil
	}

	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(fromNumber)
	params.SetBody(fmt.Sprintf("Your OTP code is: %s. It is valid for 5 minutes.", code))

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	return nil
}
