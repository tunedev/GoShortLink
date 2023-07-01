package store

import (
	"context"
	"os"
	"time"

	"github.com/tunedev/GoShortLink/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collName = os.Getenv("DB_NAME")
	dbName = os.Getenv("COLLECTION_NAME")
)

type Store struct {
	client *mongo.Client
	db *mongo.Database
}

func NewStore(uri string) (*Store, error) {
	if collName == ""{
		collName = "urls"
	}
	if dbName == ""{
		dbName = "goshortlink"
	}
	client, err := mongo.NewClient((options.Client().ApplyURI(uri)))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &Store{client: client, db: db}, nil
}

func (s *Store) SaveUrl(ctx context.Context, url *model.Url) error {
	_, err := s.db.Collection(collName).InsertOne(ctx, url)
	return err
}

func (s *Store) GetUrlByShortURL(ctx context.Context, shortUrl string) (*model.Url, error) {
	var url model.Url
	err := s.db.Collection(collName).FindOne(ctx, bson.M{"short_url": shortUrl}).Decode(&url)
	return &url, err
}
func (s *Store) GetUrlByLongURL(ctx context.Context, LongUrl string) (*model.Url, error) {
	var url model.Url
	err := s.db.Collection(collName).FindOne(ctx, bson.M{"long_url": LongUrl}).Decode(&url)
	return &url, err
}
func(s *Store) UpdateUrl(ctx context.Context, url *model.Url) error {
	_, err := s.db.Collection(collName).UpdateOne(
		ctx,
		bson.M{"_id": url.ID},
		bson.D{
			{
				"$set", bson.D{
					{"short_url", url.ShortUrl},
				},
			},
		},
	)
	return err
}