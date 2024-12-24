package app

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
)

type EmailServiceMock struct {
	apiKey string
}

func NewEmailServiceMock(apiKey string) domain.EmailSvc{
	return &EmailServiceMock{
		apiKey: apiKey,
	}
}


func (s *EmailServiceMock) SendTextEmail(ctx context.Context, receiverEmail string, tittle string, bodyText string) error{
	log.Printf("Processing email to %s", receiverEmail)
	randomWaitTimer := rand.Intn(1000)

	select{
	case <- time.After(time.Millisecond * time.Duration(randomWaitTimer)):
		log.Printf("Email to %s sent", receiverEmail)
		return nil
	case <- time.After(time.Millisecond * 500):
		log.Printf("Couldnt send email to %s", receiverEmail)
		return errors.New("")
	}
}