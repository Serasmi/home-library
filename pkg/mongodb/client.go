package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewClient(ctx context.Context, host, port, username, password, authSource string) (*mongo.Client, error) {
	var url string
	_ = authSource
	//var anonymous bool

	if username == "" || password == "" {
		//anonymous = true
		url = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		url = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(url)
	/*if !anonymous {
		opts.SetAuth(options.Credential{
			AuthSource:  authSource,
			Username:    username,
			Password:    password,
			PasswordSet: true,
		})
	}*/

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create client to mongodb due to error %w", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to create client to mongodb due to error %w", err)
	}

	return client, nil
}
