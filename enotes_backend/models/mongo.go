package models

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var once sync.Once

type Mongo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Ctx        context.Context
}

func Connect() (mongodb Mongo) {
	var cancel context.CancelFunc
	mongodb.Ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongodb.Client, err = mongo.Connect(mongodb.Ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	if err != nil {
		log.Fatalf("Unable to CONNECT MongoDB: %v", err.Error())
	}

	err = mongodb.Client.Ping(mongodb.Ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Unable to PING MongoDB: %v", err.Error())
	}

	fmt.Println("Connection to MongoDB Successful.")

	return
}

// Enter New User in mongodb. Returns an error if encountered.
func (m *Mongo) CreateNewUser(user User, db, coll string) error {

	m.Collection = m.Client.Database(db).Collection(coll)

	once.Do(func() {
		unique := true
		indexName := "email_unique_index"
		indexOptions := options.Index().SetUnique(unique).SetName(indexName)

		indexModel := mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: indexOptions,
		}

		_, err := m.Collection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatalf("Error Creating index: %v", err.Error())
		}
	})

	result, err := m.Collection.InsertOne(context.Background(), &user)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return err
	}

	fmt.Printf("result.InsertedID: %v\n", result.InsertedID)
	return nil
}

// Enter New Note in mongodb. Returns an error if encountered.
func (m *Mongo) CreateNotes(notes Notes, db, coll string) error {
	m.Collection = m.Client.Database(db).Collection(coll)

	result, err := m.Collection.InsertOne(context.Background(), &notes)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return err
	}

	fmt.Printf("result.InsertedID: %v\n", result.InsertedID)
	return nil
}

func (m *Mongo) ValidateEmailPassword(email, password string, db, coll string) (user User, err error) {
	m.Collection = m.Client.Database(db).Collection(coll)

	result := m.Collection.FindOne(context.Background(), bson.M{"email": email})
	if result.Err() != nil {
		fmt.Printf("err.Error(): %v\n", result.Err().Error())
	}

	if err = result.Decode(&user); err != nil {
		return
	}

	fmt.Printf("user: %v\n", user)

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return
	}

	fmt.Println("email - password validated!!")

	return
}

func (m *Mongo) CreateNewNote(note Notes, user primitive.ObjectID, db, coll string) error {
	m.Collection = m.Client.Database(db).Collection(coll)

	once.Do(func() {

		indexModel := mongo.IndexModel{
			Keys:    bson.M{"user": 1},
			Options: options.Index().SetName("user_index"),
		}

		_, err := m.Collection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatalf("Error Creating index on user in Notes Collection: %v", err.Error())
		}
	})

	result, err := m.Collection.InsertOne(context.Background(), &note)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return err
	}

	fmt.Printf("result.InsertedID: %v\n", result.InsertedID)
	return nil
}
