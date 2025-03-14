package account

import (
	"context"
	"fmt"
	"strconv"

	"encore.app/oauth_flows/backend/account/db"
	"encore.app/oauth_flows/backend/internal/tokens"
	"encore.dev/config"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
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
	EmailVerificationRequested, "account-email-verification-requested",
	pubsub.SubscriptionConfig[*EmailVerificationRequestedEvent]{
		Handler: handler,
	},
)

// For dependency injection.
func handler(ctx context.Context, event *EmailVerificationRequestedEvent) error {
	mailer := &GomailMailer{}
	return SendVerificationEmail(ctx, event, mailer, accountsDB)
}

func SendVerificationEmail(ctx context.Context, event *EmailVerificationRequestedEvent, mailer Mailer, database *sqldb.Database) error {
	pgxdb := sqldb.Driver(database)
	query := db.New(pgxdb)

	account, err := query.ByID(ctx, event.AccountID)
	if err != nil {
		return err
	}

	config := MailerConfig{
		SMTPHost:       mailConfig.SMTPHost(),
		SMTPPassword:   secrets.SMTPPassword,
		SMTPPort:       mailConfig.SMTPPort(),
		SMTPUsername:   mailConfig.SMTPUsername(),
		SendEmailsFrom: mailConfig.SendEmailsFrom(),
	}

	link, err := generateEmailVerificationLinkForAccount(&account)
	if err != nil {
		return err
	}

	err = mailer.SendEmail(
		account.Email,
		"Welcome to GOAuth",
		fmt.Sprintf("Welcome to GOAuth. To verify your email click this link: %s", link),
		&config,
	)
	if err != nil {
		rlog.Error("Error sending welcome email", "err", err)
		return err
	}

	if !mailConfig.SendEmails() {
		rlog.Debug("Would have sent email but it's disabled", "To", account.Email)
		return nil
	}

	return nil
}

func generateEmailVerificationLinkForAccount(account *db.Account) (string, error) {
	token, err := generateEmailVerificationTokenForAccount(account)
	if err != nil {
		return "", err
	}

	link := fmt.Sprintf("%s/verify_email?verification_token=%s", apiConfig.BaseURL(), token)

	return link, nil
}

func generateEmailVerificationTokenForAccount(account *db.Account) (string, error) {
	purpose := tokens.EmailVerification
	payload := map[string]string{"AccountID": strconv.Itoa(int(account.ID))}
	return tokens.GenerateTokenFor(purpose, payload)
}
