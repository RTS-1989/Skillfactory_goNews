package mongoDB

import (
	dbInterface "GoNews/pkg/storage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Storage Хранилище данных.
type Storage struct {
	Db *mongo.Client
}

const (
	databaseName   = "go_news" // имя учебной БД
	collectionName = "posts"   // имя коллекции в учебной БД
)

// New Конструктор, принимает строку подключения к БД.
func New(ctx context.Context, constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	s := Storage{
		Db: client,
	}
	return &s, nil
}

func (mg *Storage) Posts() ([]dbInterface.Post, error) {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var data []dbInterface.Post
	for cur.Next(context.Background()) {
		var l dbInterface.Post
		err := cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		data = append(data, l)
	}
	return data, cur.Err()
}

func (mg *Storage) AddPost(p dbInterface.Post) error {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func (mg *Storage) UpdatePost(p dbInterface.Post) error {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	filter := bson.D{{"title", p.Title}}
	update := bson.D{{"$set", bson.D{{"content", p.Content}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (mg *Storage) DeletePost(p dbInterface.Post) error {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	filter := bson.D{{"title", p.Title}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
