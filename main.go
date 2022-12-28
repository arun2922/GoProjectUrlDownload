package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"      //To work on Gin Web Framework
	guuid "github.com/google/uuid"  //To generate a unique id
	"github.com/patrickmn/go-cache" //To store web requests using go-cache
)

//allUrl is a struct that contains the Urls with their Id, Uri, SourceUri

type allUrl struct {
	Id        string `json:"Id"`
	Uri       string `json:"Uri"`
	SourceUri string `json:"SourceUri"`
}

//userUrl is a struct that contains the structure for the Urls requested by user using POST method
//It contains Uri and RetryLimit

type userUrl struct {
	Uri        string `json:"Uri"`
	RetryLimit int    `json:"RetryLimit"`
}

//allUrls is a slice with data type of allUrl which contain all the urls requested by the clients

var allUrls []allUrl

// Create a cache with a default expiration time of 24 hours, and which
// purges expired items every 24 hours
var Storage = cache.New(24*time.Hour, 24*time.Hour)

// getUrls responds with the list of all albums as JSON.

func getUrls(context *gin.Context) {

	context.IndentedJSON(http.StatusOK, allUrls) //to serialize the struct into JSON and add it to the response.
}

//addUrls adds an album from JSON received in the request body.

func addUrls(context *gin.Context) {

	var newuserUrl userUrl

	// Call BindJSON to bind the received JSON to
	// newuserUrl.
	if err := context.BindJSON(&newuserUrl); err != nil {
		return
	}

	//Check if this Uri is present in the cache or not.

	if checkinCache(newuserUrl.Uri) {
		fmt.Println("Already Downloaded")
	} else {
		var newallUrl allUrl

		newallUrl.Id = genUUID() //generate new unique id for the given post request
		newallUrl.SourceUri = newallUrl.Id + ".html"
		newallUrl.Uri = newuserUrl.Uri

		//Store the Uri into the cache with value as Id

		Storage.Set(newallUrl.Uri, newallUrl.Id, cache.DefaultExpiration)

		//append the newallUrl to allUrls

		allUrls = append(allUrls, newallUrl)

		fileUrl := newuserUrl.Uri
		retry := newuserUrl.RetryLimit

		//Next for loop tries to Download File upto minimum of 10 and RetryLimit.

		for i := 0; i <= int(math.Min(float64(10), float64(retry))); i++ {
			err := DownloadFile(newallUrl.SourceUri, fileUrl)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("Downloaded: " + fileUrl)
				break
			}
		}

		//Add a 201 status code to the response, along with JSON representing the Uri you added.
		context.IndentedJSON(http.StatusCreated, newallUrl)
	}
}
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
func genUUID() string {
	id := guuid.New()
	return id.String()
}
func checkinCache(key string) bool {
	_, found := Storage.Get(key)

	if found {
		return true
	} else {
		return false
	}
}
func main() {

	//Initialize a Gin router using Default.
	router := gin.Default()

	//Use the GET function to associate the GET HTTP method and /get path with a handler function.
	router.GET("/get", getUrls)

	router.POST("/pagesource", addUrls)

	//Use the Run function to attach the router to an http.Server and start the server.
	router.Run("localhost:7771")
}

//curl -H "Content-Type: application/json" -d "{\"Uri\":\"https\"}" -X POST http://localhost:7771/pagesource
//curl -H "Content-Type: application/json" -d "{\"Id\":\"https\",\"Uri\":\"www\",\"SourceUri\":\"hui\"}" -X POST http://localhost:7771/pagesource
