package service

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type EmailService interface {
	Send(ctx context.Context, params emailParam) error
}

type emailService struct {
	client *resend.Client
}

func NewEmailService(c *resend.Client) EmailService {
	return &emailService{
		client: c,
	}
}

func (s *emailService) Send(ctx context.Context, params emailParam) error {
	if _, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From: "",
		To:   []string{params.to},
		Html: fmt.Sprintf(params.text, params.args...),
	}); err != nil {
		return err
	}

	return nil
}

type emailParam struct {
	to   string
	text string
	args []any
}
