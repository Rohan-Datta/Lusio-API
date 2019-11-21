package main
import (
	"fmt"
	"strconv"
	//"context"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	//"google.golang.org/api/iterator"
	"time"
	"os/exec"
	"strings"
	"os"
	"io/ioutil"
)

var dirPath = "/home/rohan/Documents/Lusio/API/"
var data = dirPath+"new_data.json"

type Games struct {
    Games []Game `json:"games"`
}

type Game struct {
	Title    	 string    `json:"title"`
    Publisher    string    `json:"publisher"`
    Developer    string    `json:"developer"`
    Category	 string    `json:"category"`
    Year		 int       `json:"year"`
    Game_ID	     string    `json:"game_id"`
}

/*
func initAuth() (*firestore.Client, context.Context) {
	projectID := "lusio-30742"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
	    log.Fatalf("Failed to create client: %v", err)
	}

	return client, ctx
}
*/

func homeFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func getDetails(gameID string) interface{} {
	dataFile, err := os.Open(data)

	if err != nil {
    log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(dataFile)
	defer dataFile.Close()

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	gameDetails := result[gameID]

	return gameDetails

}

func getOneTitle(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t)
	fmt.Println("Getting one title")
	fmt.Println("----------------------------------")
	gameID := r.URL.Query().Get("game_id")

	gameDetails := getDetails(gameID)
  	js, err := json.Marshal(gameDetails)

  	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	log.Fatal(err)
  	}

	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)
}

func getAllTitles(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t)
	fmt.Println("Getting all titles")
	fmt.Println("----------------------------------")
	
	var res []interface{}
	res = append(res, getDetails("T810"), getDetails("T024"), getDetails("T811"), getDetails("T002"), getDetails("T003"), getDetails("T005"), getDetails("T006"), getDetails("T007"), getDetails("T009"), getDetails("T010"))

	js, _ := json.Marshal(res)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getMostSimilar(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t)
	fmt.Println("Getting recommended titles")
	fmt.Println("----------------------------------")

	limit := "10"
	inputID := r.URL.Query().Get("input_id")
	cmd := exec.Command("python3", dirPath+"get_most_similar.py", inputID, limit)

    out, err := cmd.CombinedOutput()
    if err != nil { fmt.Println(err); }

    resSlice := strings.Split(string(out[1:len(out)-2]), ",")

	k, _ := strconv.Atoi(limit)

	var res []interface{}
    for i := 0; i < k; i++ {
    	res = append(res, getDetails((strings.TrimSpace(resSlice[i]))[1:5]))
    }

    js, _ := json.Marshal(res)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func launchGame(w http.ResponseWriter, r *http.Request) (string, string, string) {
	userID := r.URL.Query().Get("user_id")
	gameID := r.URL.Query().Get("game_id")
	ifResume := r.URL.Query().Get("if_resume")

	return userID, gameID, ifResume
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeFunc).Methods("GET")
	router.HandleFunc("/games/", getAllTitles).Methods("GET")
	router.HandleFunc("/game/", getOneTitle).Methods("GET")
	router.HandleFunc("/recommend/", getMostSimilar).Methods("GET")
	fmt.Println("Running...")
	http.ListenAndServe(":3000", router)
}