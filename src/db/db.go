package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"todo-list/src/task"
)

var db *mongo.Client

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

func GetTasksByColumn(column task.Column) []list.Item {
	if db == nil {
		log.Fatal("You must call connect before getting collections")
	}
	var results []bson.M
	opts := options.Find().SetSort(bson.D{{"index", 1}})
	cursor, err := TaskCollection().Find(context.TODO(), bson.D{{"column", column}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	type thisTask struct {
		Index       int
		Status      task.Column
		Title       string
		Description string
	}
	var _tasks []thisTask
	err = json.Unmarshal(jsonData, &_tasks)
	if err != nil {
		fmt.Println("error:", err)
	}
	tasks := make([]list.Item, len(_tasks))

	for i, _t := range _tasks {
		tasks[i] = task.NewWithIndex(_t.Status, _t.Title, _t.Description, _t.Index)
	}

	return tasks
}
