package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Collection
}

type Instrument struct {
	Name string `json:"name" bson:"name"`
	Family string `json:"family" bson:"family"`
	Price float64 `json:"price" bson:"price"`
}

// Create Instruments
func (a *App) createInstrument(w http.ResponseWriter, r *http.Request){
	var i Instrument
	decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&i); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
	_, err := a.DB.InsertOne(context.TODO(), i)
	if err != nil {
    respondWithError(w, http.StatusBadRequest, "Error!")
	}
	respondWithJSON(w, http.StatusCreated, "Successfully created instrument!")
}

// Read All Instruments
func (a *App) getInstruments(w http.ResponseWriter, r *http.Request){
	res := []bson.M{}
	cursor, err := a.DB.Find(context.TODO(), bson.M{})
	if err != nil {
    	log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
    	var inst bson.M
    	if err = cursor.Decode(&inst); err != nil {
        	 respondWithError(w, http.StatusBadRequest, "Error!")
    	}
    	res = append(res, inst)
	}
	respondWithJSON(w, http.StatusOK, res)
}

// Read Instrument
func (a *App) getInstrument(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
    	respondWithError(w, http.StatusBadRequest, "Invalid Id!")
	}
	var res bson.M
	a.DB.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&res)
	respondWithJSON(w, http.StatusOK, res)
}

//Update Instrument
func (a *App) updateInstrument(w http.ResponseWriter, r *http.Request){
	var i Instrument
	decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&i); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
	id := mux.Vars(r)["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
    	respondWithError(w, http.StatusBadRequest, "Invalid Id!")
		return
	}
	_, err = a.DB.ReplaceOne(context.TODO(), bson.M{"_id": objectId}, i)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Could not update")
		return
	}
	respondWithJSON(w, http.StatusNoContent, i)
}

//Delete Instrument
func (a *App) deleteInstrument(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
    	respondWithError(w, http.StatusBadRequest, "Invalid Id!")
		return
	}
	_, err = a.DB.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Could not update")
		return
	}
	respondWithJSON(w, http.StatusAccepted, "Successfully Deleted Instrument")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) Initialize() {
	var err error

	// Init Router
	a.Router = mux.NewRouter()
	
	// Connect to MongoDB Instance
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
    log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
    	log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	clientDb := client.Database("instrumentsdb")
	a.DB = clientDb.Collection("instruments")

	a.Router = mux.NewRouter()
	a.initializeRouters()
}

func (a *App) initializeRouters() {
	
	// // Route Handlers / Endpoints
	a.Router.HandleFunc("/api/instruments", a.getInstruments).Methods("GET")
	a.Router.HandleFunc("/api/instruments/{id}", a.getInstrument).Methods("GET")
	a.Router.HandleFunc("/api/instruments", a.createInstrument).Methods("POST")
	a.Router.HandleFunc("/api/instruments/{id}", a.updateInstrument).Methods("PUT")
	a.Router.HandleFunc("/api/instruments/{id}", a.deleteInstrument).Methods("DELETE")
}

func (a *App) Run(addr string) {
	fmt.Println("Server Listening on Port 8000")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}