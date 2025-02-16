package user

import "encore.dev/pubsub"

// type SignupEvent struct{ UserID int }
//
// var Signups = pubsub.NewTopic[*SignupEvent]("signups", pubsub.TopicConfig{
// 	DeliveryGuarantee: pubsub.AtLeastOnce,
// })

type EmailVerificationRequestedEvent struct{ UserID int }

var EmailVerificationRequested = pubsub.NewTopic[*EmailVerificationRequestedEvent]("email-verification", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
