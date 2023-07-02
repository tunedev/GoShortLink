package store

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/tunedev/GoShortLink/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getIPNodeID() (int64, error){
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ip := localAddr.IP.To4()
	lastSegment := ip[len(ip)-1]
	return int64(lastSegment), nil
}

var (
	dbName = os.Getenv("DB_NAME")
	collName = os.Getenv("COLLECTION_NAME")
)

type Store struct {
	client *mongo.Client
	db *mongo.Database
	node *snowflake.Node
}

func NewStore(uri string) (*Store, error) {
	if collName == ""{
		collName = "urls"
	}
	if dbName == ""{
		dbName = "goshortlink"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, (options.Client().ApplyURI(uri)))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	nodeID, err := getIPNodeID()
	if err != nil {
		return nil, err
	}
	
	// Create a new Snowflake node with the node ID
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("could not create Snowflake node: %w", err)
	}
	
	return &Store{client: client, db: db, node: node}, nil
}

func (s *Store) SaveUrl(ctx context.Context, url *model.Url) error {
	id := s.node.Generate()

	url.ID = "id_" + id.String()
	_, err := s.db.Collection(collName).InsertOne(ctx, url)
    if err != nil {
        return err
    }

    return nil
}

func (s *Store) GetUrlByShortURL(ctx context.Context, shortUrl string) (*model.Url, error) {
	var url model.Url
	err := s.db.Collection(collName).FindOne(ctx, bson.M{"short_url": shortUrl}).Decode(&url)
	return &url, err
}

func (s *Store) GetUrlByLongURL(ctx context.Context, LongUrl string) (*model.Url, error) {
	var url model.Url
    fmt.Printf("Finding URL: %s\n", LongUrl)
    err := s.db.Collection(collName).FindOne(ctx, bson.M{"long_url": LongUrl}).Decode(&url)
    if err != nil {
        fmt.Printf("Error finding URL: %v\n", err)
        return nil, err
    }
    return &url, err
}

func(s *Store) UpdateUrl(ctx context.Context, url *model.Url) error {
    filter := bson.M{"_id": url.ID} // no conversion here
    update := bson.M{
        "$set": bson.M{
            "short_url": url.ShortUrl,
        },
    }
    result, err := s.db.Collection(collName).UpdateOne(ctx, filter, update)
    if err != nil {
        fmt.Printf("Error updating URL: %v\n", err)
        return err
    }
    fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
    return nil
}

func (s *Store) GetAllUrls(ctx context.Context) ([]*model.Url, error) {
	var urls []*model.Url

	cursor, err := s.db.Collection(collName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var url model.Url
		if err := cursor.Decode(&url); err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}