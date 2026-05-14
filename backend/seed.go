//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Question struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	TestID  string             `bson:"test_id"`
	Text    string             `bson:"text"`
	Options []string           `bson:"options"`
	Answer  int                `bson:"answer"`
	Subject string             `bson:"subject"`
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database("rayan_academy").Collection("questions")

	questions := []interface{}{
		Question{TestID: "ssc-cgl-mock-1", Text: "What is the capital of India?", Options: []string{"Delhi", "Mumbai", "Kolkata", "Chennai"}, Answer: 0, Subject: "General Knowledge"},
		Question{TestID: "ssc-cgl-mock-1", Text: "Who wrote the national anthem of India?", Options: []string{"Rabindranath Tagore", "Mahatma Gandhi", "Jawaharlal Nehru", "B. R. Ambedkar"}, Answer: 0, Subject: "General Knowledge"},
		Question{TestID: "ssc-cgl-mock-1", Text: "What is the currency of Japan?", Options: []string{"Yuan", "Won", "Yen", "Ringgit"}, Answer: 2, Subject: "General Knowledge"},
		Question{TestID: "bank-po-1", Text: "If x+5=10, what is x?", Options: []string{"3", "4", "5", "6"}, Answer: 2, Subject: "Quantitative Aptitude"},
		Question{TestID: "bank-po-1", Text: "What is 15% of 200?", Options: []string{"20", "25", "30", "35"}, Answer: 2, Subject: "Quantitative Aptitude"},
		Question{TestID: "bank-po-1", Text: "A train travels 120 km in 2 hours. What is its speed in km/h?", Options: []string{"50", "60", "70", "80"}, Answer: 1, Subject: "Quantitative Aptitude"},
		Question{TestID: "ssc-cgl-mock-1", Text: "Which planet is known as the Red Planet?", Options: []string{"Earth", "Mars", "Jupiter", "Venus"}, Answer: 1, Subject: "General Science"},
		Question{TestID: "ssc-cgl-mock-1", Text: "What is the chemical symbol for water?", Options: []string{"H2O", "CO2", "O2", "NaCl"}, Answer: 0, Subject: "General Science"},
		Question{TestID: "bank-po-1", Text: "Find the missing number: 2,4,8,16,?", Options: []string{"24", "32", "30", "28"}, Answer: 1, Subject: "Reasoning"},
		Question{TestID: "bank-po-1", Text: "If CAT is coded as 3120, how is DOG coded?", Options: []string{"4157", "5147", "4158", "5148"}, Answer: 0, Subject: "Reasoning"},
	}

	_, err = coll.InsertMany(context.TODO(), questions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ 10 seed questions inserted!")
}
