package account

import "encore.dev/pubsub"

type EmailVerificationRequestedEvent struct{ UserID int64 }

var EmailVerificationRequested = pubsub.NewTopic[*EmailVerificationRequestedEvent]("account-email-verification", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
