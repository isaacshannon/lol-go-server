package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User stores data about a user
type User struct {
	UserID string
	MapX0  float64
	MapX1  float64
	MapY0  float64
	MapY1  float64
}

func main() {
	ctx := context.TODO()

	mongoURI, ok := os.LookupEnv("DB_URI")
	if !ok {
		fmt.Println("error: unable to find MONGO_PW in the environment")
		os.Exit(1)
	}
	fmt.Println("connection string is:", mongoURI)

	// Set client options and connect
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println("Connected to MongoDB!")
	//
	//collection := client.Database("league").Collection("users")
	//ash := User{"Ash", 10., 11.,12., 13.}
	//insertResult, err := collection.InsertOne(context.TODO(), ash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	//
	//filter := bson.D{}
	//var result User
	//err = collection.FindOne(context.TODO(), filter).Decode(&result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print(result)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/predict", servePredict)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

type predictResponse struct {
	Result string `json:"result"`
}
func servePredict(w http.ResponseWriter, r *http.Request) {
	img := r.FormValue("imgBase64")
	usr := r.FormValue("user")
	resp, err := http.PostForm(
		"http://league-nodeport-service/predict",
		url.Values{"imgBase64": {img}, "user": {usr}})
	if err != nil {
		data := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		data := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
		return
	}

	pred := predictResponse{}
	err = json.Unmarshal(body, &pred)
	if err != nil {
		data := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
		return
	}

	//resp, err := http.Get("http://home-nodeport-service/predict")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//data := map[string]string{"result": img}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pred)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", "capture.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

}
