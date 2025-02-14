package user

import "encore.dev/pubsub"

type SignupEvent struct{ UserID int }

var Signups = pubsub.NewTopic[*SignupEvent]("signups", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
