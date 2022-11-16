package db

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"todo-list/src/task"
)

var db *mongo.Client

type MyObjectID string

func (id MyObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	p, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bsontype.Null, nil, err
	}

	return bson.MarshalValue(p)
}

func Connect() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db = client
	return client
}

func TaskCollection() *mongo.Collection {
	if db == nil {
		log.Fatal("You must call connect before getting collections")
	}
	return db.Database("app").Collection("tasks")
}

func GetTasksByColumn(status task.Column) []list.Item {
	if db == nil {
		log.Fatal("You must call connect before getting collections")
	}
	type thisTask struct {
		ID          string `bson:"_id"`
		Index       int
		Status      task.Column
		Title       string
		Description string
	}
	var _tasks []thisTask
	opts := options.Find().SetSort(bson.D{{"index", 1}})
	cursor, err := TaskCollection().Find(context.TODO(), bson.D{{"status", status}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(context.TODO(), &_tasks)
	if err != nil {
		log.Fatal(err)
	}

	tasks := make([]list.Item, len(_tasks))

	for i, _t := range _tasks {
		_id, err := primitive.ObjectIDFromHex(_t.ID)
		if err != nil {
			fmt.Println("Err:", err)
		}
		tasks[i] = task.NewWithIndex(_t.Status, _t.Title, _t.Description, _t.Index, _id)
	}

	return tasks
}

func DeleteTaskByID(ID primitive.ObjectID) {
	filter := bson.D{{"_id", ID}}
	_, err := TaskCollection().DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}

func InsertTask(index int, t task.Task) {
	doc := bson.D{
		{"index", index},
		{"status", t.Status},
		{"title", t.Title()},
		{"description", t.Description()},
	}
	_, err := TaskCollection().InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}
