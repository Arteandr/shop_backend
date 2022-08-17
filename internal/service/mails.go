package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"shop_backend/internal/repository"
	apperrors "shop_backend/pkg/errors"
	"shop_backend/pkg/mail"
)

type MailsService struct {
	repo       repository.Mails
	mailSender mail.Sender
}

func NewMailsService(repo repository.Mails, mailSender mail.Sender) *MailsService {
	return &MailsService{
		repo:       repo,
		mailSender: mailSender,
	}
}

func (s *MailsService) CreateVerify(ctx context.Context, userId int, login, email string) error {
	token := uuid.New().String()
	err := s.mailSender.SendVerify(email, login, token)
	if err != nil {
		fmt.Println(err)
		return apperrors.ErrEmailSend
	}

	if err := s.repo.SetVerify(ctx, token, userId); err != nil {
		return err
	}

	return nil
}
