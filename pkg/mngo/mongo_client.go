package mngo

import (
	"context"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnectionURL(
	credentials BasicAuthCredentials,
	host string,
	database string,
) string {
	u := &url.URL{
		Scheme: "mongodb+srv",
		User:   url.UserPassword(credentials.Username, credentials.Password),
		Host:   host,
		Path:   database,
		RawQuery: url.Values{
			"retryWrites": []string{"true"},
			"w":           []string{"majority"},
		}.Encode(),
	}

	return u.String()
}

type BasicAuthCredentials struct {
	Username string `json:"-"`
	Password string `json:"-"`
}

func NewMongoDatabase(
	ctx context.Context,
	credentials BasicAuthCredentials,
	host string,
	database string,
) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(NewMongoConnectionURL(credentials, host, database)))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(database), nil
}
