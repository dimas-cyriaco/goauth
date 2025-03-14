package account

import "encore.dev/pubsub"

type EmailVerificationRequestedEvent struct{ AccountID int64 }

var EmailVerificationRequested = pubsub.NewTopic[*EmailVerificationRequestedEvent]("account-email-verification", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
