package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/zillalikestocode/community-api/configs"
	"github.com/zillalikestocode/community-api/models"
	"github.com/zillalikestocode/community-api/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Community struct {
}

var communityCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "communities")

// get user communities
func (c *Community) GetAll(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var result []bson.M

	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))

	cursor, _ := communityCollection.Find(context.TODO(), bson.M{"members": bson.M{"$elemMatch": bson.M{"id": userId}}})

	if err := cursor.All(context.TODO(), &result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "A server error has occured", Data: map[string]interface{}{"error": err.Error()}})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusOK, Message: "Communities fetched successfully", Data: map[string]interface{}{"result": result}})
}

// create community
func (c *Community) Create(w http.ResponseWriter, r *http.Request) {
	var body models.Community
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "Please pass the required details", Data: map[string]interface{}{"error": err.Error()}})
		return
	}

	newCommunity := models.Community{
		ID:          primitive.NewObjectID(),
		Name:        body.Name,
		Description: body.Description,
		Owner:       userId,
		Members:     []models.Member{{ID: userId, Admin: true}},
	}
	result, err := communityCollection.InsertOne(context.TODO(), newCommunity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "An error occured while creating the community", Data: map[string]interface{}{"error": err.Error()}})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusCreated, Message: "Community created", Data: map[string]interface{}{"community": result}})
}

// join community
func (c *Community) Join(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CommunityId string `json:"communityId"`
	}
	var community models.Community
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "Please pass the required details", Data: map[string]interface{}{"error": err.Error()}})
		return
	}
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)
	if err := communityCollection.FindOne(context.TODO(), bson.M{"_id": communityId}).Decode(&community); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "Unable to find community", Data: map[string]interface{}{"error": err.Error(), "id": communityId}})
		return
	}
	communityMembers := community.Members
	if slices.Contains(communityMembers, models.Member{ID: userId, Admin: false}) || slices.Contains(communityMembers, models.Member{ID: userId, Admin: true}) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "User already in the community"})
		return
	} else {
		result, err := communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId}, bson.M{"$push": bson.M{"members": userId}})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "Couldn't join community", Data: map[string]interface{}{"error": err.Error()}})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusOK, Message: "Successfully joined community", Data: map[string]interface{}{"result": result}})

	}
}

// leave a community
func (c *Community) Leave(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var body struct {
		CommunityId string `json:"communityId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadGateway, Message: "Please pass in the required body parameters", Data: map[string]interface{}{"error": err.Error()}})
		return
	}
	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)
	result, err := communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId}, bson.M{"$pull": bson.M{"members": bson.M{"id": userId}}})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadGateway, Message: "An error occured while leaving the community", Data: map[string]interface{}{"error": err.Error()}})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusOK, Message: "Successfully left the community", Data: map[string]interface{}{"result": result}})

	}
}

// search community
func (c *Community) SearchCommunity(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	var result []bson.M

	cursor, _ := communityCollection.Find(context.TODO(), bson.M{"name": bson.M{"$regex": query}})

	if err := cursor.All(context.TODO(), &result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadRequest, Message: "A server error has occured"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusOK, Message: "Communities found", Data: map[string]interface{}{"result": result}})

}

// ANNOUNCEMENT SECTION

// create announcement
func (c *Community) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var body struct {
		Name        string `json:"name"`
		Date        string `json:"date"`
		Message     string `json:"message"`
		CommunityId string `json:"communityId"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	parsedDate, _ := time.Parse(time.RFC3339, body.Date)
	date := primitive.NewDateTimeFromTime(parsedDate)

	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)

	newAnnouncement := bson.M{"announcements": bson.M{"date": date, "message": body.Message, "id": primitive.NewObjectID(), "creator": bson.M{"name": body.Name, "id": userId}}}

	result, err := communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId}, bson.M{"$push": newAnnouncement})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadGateway, Message: "An error occured while creating an announcement", Data: map[string]interface{}{"error": err.Error()}})
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusCreated, Message: "Announcement created successfully", Data: map[string]interface{}{"result": result}})
	}

}

// delete announcement
func (c *Community) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var body struct {
		AnnouncementId string `json:"announcementId"`
		CommunityId    string `json:"communityId"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)
	announcementId, _ := primitive.ObjectIDFromHex(body.AnnouncementId)

	communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId}, bson.M{"$pull": bson.M{"announcements": bson.M{"id": announcementId, "creator.id": userId}}})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusOK, Message: "Announcement deleted"})
}

// EVENTS SECTIONS

// create event
func (c *Community) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// _, claims, _ := jwtauth.FromContext(r.Context())
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Time        string `json:"time"`
		CommunityId string `json:"communityId"`
		Address     string `json:"address"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	// userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)
	parsedDate, _ := time.Parse(time.RFC3339, body.Date)
	date := primitive.NewDateTimeFromTime(parsedDate)

	newEvent := bson.M{
		"name":        body.Name,
		"id":          primitive.NewObjectID(),
		"description": body.Description,
		"date":        date,
		"time":        body.Time,
		"address":     body.Address,
	}

	result, err := communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId}, bson.M{"$push": bson.M{"events": newEvent}})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadGateway, Message: "An error occured while adding the event", Data: map[string]interface{}{"error": err.Error()}})
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusCreated, Message: "Event added successfully", Data: map[string]interface{}{"result": result}})
	}

}

// delete event
func (c *Community) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var body struct {
		EventId     string `json:"eventId"`
		CommunityId string `json:"communityId"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)
	eventId, _ := primitive.ObjectIDFromHex(body.EventId)

	communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId, "owner": userId}, bson.M{"$pull": bson.M{"events": bson.M{"id": eventId}}})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusOK, Message: "Event deleted"})
}

// update event
func (c *Community) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Time        string `json:"time"`
		CommunityId string `json:"communityId"`
		EventId     string `json:"eventId"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	// userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	communityId, _ := primitive.ObjectIDFromHex(body.CommunityId)
	eventId, _ := primitive.ObjectIDFromHex(body.EventId)
	parsedDate, _ := time.Parse(time.RFC3339, body.Date)
	date := primitive.NewDateTimeFromTime(parsedDate)

	result, err := communityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId, "events.id": eventId}, bson.M{"$set": bson.M{"events.$.name": body.Name, "events.$.description": body.Description, "events.$.date": date, "events.$.time": body.Time}})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusBadGateway, Message: "An error occured while updating the event", Data: map[string]interface{}{"error": err.Error()}})
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responses.UserResponse{Status: http.StatusCreated, Message: "Event updated successfully", Data: map[string]interface{}{"result": result}})
	}

}
