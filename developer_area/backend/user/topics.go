package user

import "encore.dev/pubsub"

type EmailVerificationRequestedEvent struct{ UserID int }

var EmailVerificationRequested = pubsub.NewTopic[*EmailVerificationRequestedEvent]("email-verification", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
