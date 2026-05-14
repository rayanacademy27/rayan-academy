package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Question represents a single test question
type Question struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TestID  string             `bson:"test_id" json:"test_id"`
	Text    string             `bson:"text" json:"text"`
	Options []string           `bson:"options" json:"options"`
	Answer  int                `bson:"answer" json:"-"` // index 0-3; "-" means never send to frontend
	Subject string             `bson:"subject" json:"subject"`
}

// User is the shape of a person using our platform
type User struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

// JWT secret key – in real life, use something long and random (keep it safe!)
var jwtKey = []byte("rayan_academy_secret_key_change_me_later")

var client *mongo.Client

func signupHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	// Read the JSON body into the user variable
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if email is already taken
	usersColl := client.Database("rayan_academy").Collection("users")
	existing := usersColl.FindOne(r.Context(), map[string]string{"email": user.Email})
	if existing.Err() == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash the password (we NEVER store it as plain text)
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashed)

	// Save the user to MongoDB
	_, err = usersColl.InsertOne(r.Context(), user)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Success!
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Find the user by email
	usersColl := client.Database("rayan_academy").Collection("users")
	var stored User
	err := usersColl.FindOne(r.Context(), map[string]string{"email": creds.Email}).Decode(&stored)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the given password with the hashed one
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Create a JWT token that expires in 72 hours
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": creds.Email,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Send the token back
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// getTestQuestionsHandler returns questions for a given test ID (without answers)
func getTestQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// URL: /tests/{testId}/questions
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] == "" || parts[3] != "questions" {
		http.NotFound(w, r)
		return
	}
	testID := parts[2]

	coll := client.Database("rayan_academy").Collection("questions")
	cursor, err := coll.Find(r.Context(), bson.M{"test_id": testID})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var questions []Question
	if err = cursor.All(r.Context(), &questions); err != nil {
		http.Error(w, "Error reading questions", http.StatusInternalServerError)
		return
	}

	// Because of json:"-" tag, the Answer field is automatically omitted
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// submitTestHandler accepts user answers, calculates score, and saves the attempt
func submitTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	// URL: /tests/{testId}/submit
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] == "" || parts[3] != "submit" {
		http.NotFound(w, r)
		return
	}
	testID := parts[2]

	// Parse the incoming answer list
	var submission struct {
		Answers []struct {
			QuestionID string `json:"questionId"`
			Selected   int    `json:"selected"`
		} `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Fetch the actual questions to get correct answers
	coll := client.Database("rayan_academy").Collection("questions")
	var questions []Question
	cursor, err := coll.Find(r.Context(), bson.M{"test_id": testID})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if err = cursor.All(r.Context(), &questions); err != nil {
		http.Error(w, "Error reading questions", http.StatusInternalServerError)
		return
	}

	// Build a fast map of questionID -> correct answer
	answerMap := make(map[string]int)
	for _, q := range questions {
		answerMap[q.ID.Hex()] = q.Answer
	}

	// Calculate score
	score := 0
	total := len(questions)
	for _, ans := range submission.Answers {
		if correct, exists := answerMap[ans.QuestionID]; exists && ans.Selected == correct {
			score++
		}
	}
	percentage := float64(score) / float64(total) * 100

	// (Later we'll save the attempt to the user's history)

	result := map[string]interface{}{
		"score":      score,
		"total":      total,
		"percentage": percentage,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("✅ Connected to MongoDB!")

	// 1. Set up your Routes (using mux)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/signup", signupHandler)
	mux.HandleFunc("/login", loginHandler)

	// 2. Add the /tests/ route (with sub-routing logic)
	// I have wrapped this in authMiddleware so only logged-in users can take tests
	mux.HandleFunc("/tests/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Simple sub-router:
		if strings.HasSuffix(r.URL.Path, "/questions") && r.Method == http.MethodGet {
			getTestQuestionsHandler(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/submit") && r.Method == http.MethodPost {
			submitTestHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))

	// 3. Set up CORS Configuration
	// This allows your React frontend (localhost:5173) to talk to this backend
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// 4. Wrap the mux with the CORS handler
	handler := corsHandler.Handler(mux)

	// 5. Start server
	log.Println("🚀 Rayan Academy backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "ok",
		"message": "Rayan Academy Backend is Alive",
		"time":    time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}
