package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/zillalikestocode/community-api/app/configs"
	"github.com/zillalikestocode/community-api/app/models"
	"github.com/zillalikestocode/community-api/app/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
}

var userCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "users")

// user account creation handler
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "An error has occured",
			Data:    map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Name:     user.Name,
		Password: string(hash),
		Email:    user.Email,
	}

	result, err := userCollection.InsertOne(context.TODO(), newUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Unable to create user",
			Data:    map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    map[string]interface{}{"data": result},
	}
	json.NewEncoder(w).Encode(response)
}

// user login handler
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	tokenAuth := configs.UseJWT()

	fmt.Println("User login endpoint called")

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return
	}

	err := userCollection.FindOne(context.TODO(), bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "User does not exist", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "Incorrect Password", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	claims := map[string]interface{}{"id": user.ID, "email": user.Email}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*336))
	_, tokenString, _ := tokenAuth.Encode(claims)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Log in Successfull",
		Data:    map[string]interface{}{"token": tokenString},
	}
	json.NewEncoder(w).Encode(response)
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtbWFudWVsbmdva2E3NzhAZ21haWwuY29tIiwiaWQiOiI2NWRiYjZlOTBiZWU4ZWQ1NjdjMGViNjkifQ.D2ya5jgOKJzo6bRVZKNiEewfsx23stKZmgjFWRkfbtI
// get user with token
func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var user models.User
	objectId, _ := primitive.ObjectIDFromHex(claims["id"].(string))

	if err := userCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "An error has occurred",
			Data:    map[string]interface{}{"error": err.Error(), "id": claims["id"]},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	response := responses.UserResponse{
		Status:  http.StatusOK,
		Message: "User successfully fetched",
		Data:    map[string]interface{}{"user": user},
	}
	json.NewEncoder(w).Encode(response)

}
func (u *User) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User deletion endpoint called")
}
func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User update endpoint called")
}
