package mock

import "context"

type EmailMock struct{}

func (e *EmailMock) SendTextEmail(ctx context.Context, receiverEmail string, tittle string, bodyText string) error {
	return nil
}