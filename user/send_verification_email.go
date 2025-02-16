package user

import (
	"context"
	"fmt"
	"strconv"

	tokengenerator "encore.app/internal/token_generator"
	"encore.app/utils"
	"encore.dev/config"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MailConfig struct {
	SendEmails     config.Bool   // Whether or not to send emails on current environment
	SendEmailsFrom config.String // Email address to be used as sender
	SMTPHost       config.String // SMTP server host
	SMTPPort       config.Int    // SMTP server port
	SMTPUsername   config.String // SMTP server username
}

var mailConfig *MailConfig = config.Load[*MailConfig]()

type APIConfig struct {
	BaseURL config.String // API Base URL
}

var apiConfig *APIConfig = config.Load[*APIConfig]()

var secrets struct {
	SMTPPassword string // SMTP server password
}

var _ = pubsub.NewSubscription(
	EmailVerificationRequested, "email-verification-requested",
	pubsub.SubscriptionConfig[*EmailVerificationRequestedEvent]{
		Handler: handler,
	},
)

// For dependency injection.
func handler(ctx context.Context, event *EmailVerificationRequestedEvent) error {
	mailer := &utils.GomailMailer{}
	return SendVerificationEmail(ctx, event, mailer)
}

func SendVerificationEmail(ctx context.Context, event *EmailVerificationRequestedEvent, mailer utils.Mailer) error {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: db.Stdlib()}))
	if err != nil {
		return err
	}

	var user User
	if err = db.Where("id = $1", event.UserID).First(&user).Error; err != nil {
		return err
	}

	config := utils.MailerConfig{
		SMTPHost:       mailConfig.SMTPHost(),
		SMTPPassword:   secrets.SMTPPassword,
		SMTPPort:       mailConfig.SMTPPort(),
		SMTPUsername:   mailConfig.SMTPUsername(),
		SendEmailsFrom: mailConfig.SendEmailsFrom(),
	}

	link, err := generateEmailVerificationLinkForUser(&user)
	if err != nil {
		return err
	}

	err = mailer.SendEmail(
		user.Email,
		"Welcome to GOAuth",
		fmt.Sprintf("Welcome to GOAuth. To verify your email click this link: %s", link),
		&config,
	)
	if err != nil {
		rlog.Error("Error sending welcome email", "err", err)
		return err
	}

	if !mailConfig.SendEmails() {
		rlog.Debug("Would have sent email but it's disabled", "To", user.Email)
		return nil
	}

	return nil
}

func generateEmailVerificationLinkForUser(user *User) (string, error) {
	token, err := generateEmailVerificationTokenForUser(user)
	if err != nil {
		return "", err
	}

	link := fmt.Sprintf("%s/verify_email?verification_token=%s", apiConfig.BaseURL(), token)

	return link, nil
}

func generateEmailVerificationTokenForUser(user *User) (string, error) {
	purpose := tokengenerator.EmailVerification
	payload := map[string]string{"UserID": strconv.Itoa(int(user.ID))}
	return tokengenerator.GenerateTokenFor(purpose, payload)
}
