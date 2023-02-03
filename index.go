package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type DefaultBangalore struct {
	Lat string
	Lng string
}

type IconMarker struct {
	UrlImgMarker string
	Size         []string
	Animation    string
	Draggable    string
	Numbering    string
}

type DataMap struct {
	MapId            string
	Zoom             string
	LogoUrl          string
	DisableDefaultUI string
	AllMarkers       string
	DefaultBangalore
	IconMarker
}

type Page struct {
	Href         string `json:"href" bson:"href"`
	UrlImgMarker string `json:"urlImgMarker" bson:"urlImgMarker"`
	Bangalore    `json:"bangalore" bson:"bangalore"`
	DataPopup    `json:"dataPopup" bson:"dataPopup"`
}

type DataPopup struct {
	Title  string `json:"title" bson:"title"`
	Text   string `json:"text" bson:"text"`
	UrlImg string `json:"urlImg" bson:"urlImg"`
	Links  []Link `json:"links" bson:"links"`
}
type Link struct {
	Url  string `json:"url" bson:"url"`
	Name string `json:"name" bson:"name"`
}
type Bangalore struct {
	Lat int `json:"lat" bson:"lat"`
	Lng int `json:"lng" bson:"lng"`
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.LoadHTMLGlob("static/index.html")
	api := router.Group("/api")
	{
		api.POST("/post-general-data", postGeneralData)
		api.POST("/post-pages-data", postPagesData)
		api.GET("/get-general-data", getGeneralData)
		api.GET("/get-pages-data", getPagesData)
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.Run("localhost:5050")

}

func getGeneralData(c *gin.Context) {

}

func getPagesData(c *gin.Context) {
	var allPages = []Page{}
	var result Page
	db, err := newClient(context.TODO(), "boss", "amurasila", "googleMap")
	if err != nil {
		log.Fatalf("failed to connect to mongo db")
	}

	collection := db.Collection("dataPages")

	cur, err := collection.Find(context.Background(), bson.D{})
	for cur.Next(context.Background()) {
		err = cur.Decode(&result)
		if err != nil {
			log.Fatalf("failed to decode data page. error: %v", err)
		}
		allPages = append(allPages, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatalf("failed error: %v", err)
	}
	cur.Close(context.TODO())

	c.IndentedJSON(http.StatusOK, allPages)
	fmt.Println(result)

}

func postGeneralData(c *gin.Context) {
	var data DataMap

	if err := c.BindJSON(&data); err != nil {
		return
	}

	fmt.Println(data)
	c.IndentedJSON(http.StatusCreated, data)
}

func postPagesData(c *gin.Context) {
	var page Page

	if err := c.BindJSON(&page); err != nil {
		return
	}
	db, err := newClient(context.TODO(), "boss", "amurasila", "googleMap")
	if err != nil {
		log.Fatalf("error to connect to mongo db. error %v", err)
	}

	collection := db.Collection("dataPages")
	res, err := collection.InsertOne(context.TODO(), page)
	if err != nil {
		log.Fatalf("failed to create page. error: %v", err)
	}
	fmt.Println(res)
	fmt.Println(page)
	c.IndentedJSON(http.StatusCreated, page)
}

func newClient(ctx context.Context, username, password, database string) (*mongo.Database, error) {
	mongoDBURL := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.mgdx7q8.mongodb.net/?retryWrites=true&w=majority", username, password)
	clientOptions := options.Client().ApplyURI(mongoDBURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed connect to mongodb: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("failed to ping to mongodb: %v", err)
	}
	return client.Database(database), err
}
