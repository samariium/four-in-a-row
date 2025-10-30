package store

import (
	"context"
	"time"
	"fmt"   
	"github.com/yourname/fourinarow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Client     *mongo.Client
	DB         *mongo.Database
	PlayersCol *mongo.Collection
	GamesCol   *mongo.Collection
}

func NewMongoStore(ctx context.Context, uri string) (*MongoStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil { return nil, err }
	if err := client.Connect(ctx); err != nil { return nil, err }
	db := client.Database("fourinarow")
	return &MongoStore{
		Client:     client,
		DB:         db,
		PlayersCol: db.Collection("players"),
		GamesCol:   db.Collection("games"),
	}, nil
}

func (s *MongoStore) EnsurePlayer(ctx context.Context, username string) error {
	_, err := s.PlayersCol.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$setOnInsert": bson.M{
			"username": username, "wins": 0, "losses": 0, "draws": 0,
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (s *MongoStore) IncWinLoss(ctx context.Context, winner, loser string) error {
	if winner != "BOT" {
		_, _ = s.PlayersCol.UpdateOne(ctx, bson.M{"username": winner}, bson.M{"$inc": bson.M{"wins": 1}})
	}
	if loser != "BOT" {
		_, _ = s.PlayersCol.UpdateOne(ctx, bson.M{"username": loser}, bson.M{"$inc": bson.M{"losses": 1}})
	}
	return nil
}

func (s *MongoStore) IncDraws(ctx context.Context, users []string) error {
	_, err := s.PlayersCol.UpdateMany(ctx,
		bson.M{"username": bson.M{"$in": users}},
		bson.M{"$inc": bson.M{"draws": 1}},
	)
	return err
}

func (s *MongoStore) InsertGame(ctx context.Context, g models.GameDoc) error {
	g.CreatedAt = time.Now()
	_, err := s.GamesCol.InsertOne(ctx, g)
	return err
}

func (s *MongoStore) TopPlayers(ctx context.Context, limit int64) ([]models.Player, error) {
    opts := options.Find().SetSort(bson.D{
        {Key: "wins",  Value: -1},
        {Key: "draws", Value: -1},
    }).SetLimit(limit)

    cur, err := s.PlayersCol.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, fmt.Errorf("find players: %w", err)
    }
    defer cur.Close(ctx)

    var out []models.Player
    for cur.Next(ctx) {
        var p models.Player
        if err := cur.Decode(&p); err != nil {
            return nil, fmt.Errorf("decode player: %w", err)
        }
        out = append(out, p)
    }
    if err := cur.Err(); err != nil {
        return nil, fmt.Errorf("cursor: %w", err)
    }
    return out, nil
}

