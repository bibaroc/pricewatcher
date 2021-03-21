package msg

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	mmr                 struct{ collection *mongo.Collection }
	mongoMessageWrapper struct {
		ID string `json:"_id,omitempty" bson:"_id"`
		Message
	}
)

func (r *mmr) Save(ctx context.Context, id string, text string) error {
	_, err := r.collection.InsertOne(ctx,
		mongoMessageWrapper{
			ID: id,
			Message: Message{
				CreatedAt: time.Now().UTC(),
				Text:      text,
			}},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *mmr) Get(ctx context.Context, id string) (*Message, error) {
	message := mongoMessageWrapper{}
	if err := r.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&message); err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &message.Message, nil
}

func NewMongoMessageRepository(db *mongo.Database) (MessageRepository, error) {
	return &mmr{db.Collection("messages")}, nil
}
