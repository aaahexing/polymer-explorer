package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	CategoryResultLabelInKG          = "kg"
	CategoryResultLabelInTrainingset = "training_set"
	CategoryResultLabelInProduction  = "production"
)

/*
   Data definitions, will be moved to the *proto file.
*/

// Type assertions with the corresponding type scores.
type CategoryResult struct {
	// Name of the category.
	Name string
	// Score of the category.
	Score float64
	// Labels of this category, the values can be:
	//  "kg", "training", "production".
	Labels []string
}

// The response of the entity server.
type EntityServerResponse struct {
	// The mid specified in the request.
	Mid string
	// The build specified in the request.
	Path string
	// The impression of this mid
	Impression float64
	// The assertions on this mid by the specified wpcat-build.
	Results []CategoryResult
}

/*
   HTTPUtils, will be moved to the http_utils.go.
*/

/*
   Handlers, will be moved to the handlers.go.
*/
func SendJSON(w http.ResponseWriter, v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("Failed to marshal data into json: %v", err),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// '/buildList/'
func getBuildList() ([]string, error) {
	tempBuildList := []string{
		"/cns/pc-d/home/wikiquality-prod/wpcat/daily/2015-06-08",
		"/cns/pc-d/home/wikiquality-prod/wpcat/daily/2015-06-10"}
	return tempBuildList, nil
}

func serveBuildListHandler(w http.ResponseWriter, r *http.Request) {
	buildList, err := getBuildList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SendJSON(w, buildList)
}

// '/entityScores/'
func GetCategoryResults(buildPath string, mid string) ([]CategoryResult, error) {
	tempCategoryResult := []CategoryResult{
		{
			Name:  "/book/author",
			Score: 0.98989898,
			Labels: []string{
				CategoryResultLabelInKG,
				CategoryResultLabelInProduction},
		},
		{
			Name:  "/people/person",
			Score: 0.98888888,
			Labels: []string{
				CategoryResultLabelInKG,
				CategoryResultLabelInTrainingset,
				CategoryResultLabelInProduction},
		},
		{
			Name:   "/music/artist",
			Score:  0.97786668,
			Labels: []string{},
		},
		{
			Name:  "/tv/tv_actor",
			Score: 0.878787878,
			Labels: []string{
				CategoryResultLabelInKG},
		},
		{
			Name:  "/location/location",
			Score: 0.696666666,
			Labels: []string{
				CategoryResultLabelInKG},
		},
		{
			Name:   "/film/film",
			Score:  0.543333333,
			Labels: []string{},
		},
	}
	return tempCategoryResult, nil
}

func serveEntityScoresHandler(w http.ResponseWriter, r *http.Request) {
	response := EntityServerResponse{
		Mid:  strings.TrimSpace(r.FormValue("mid")),
		Path: strings.TrimSpace(r.FormValue("path")),
	}
	category_results, err := GetCategoryResults(response.Path, response.Mid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Results = category_results
	fmt.Printf("Response of entity server: %v\n", response)
	SendJSON(w, response)
}

func serveIndexHandler(w http.ResponseWriter, r *http.Request) {
	indexFilepath, err := filepath.Abs("./static/index.html")
	if err != nil {
		fmt.Printf("Error occurs when generate absolute path: %v\n", err.Error)
	}
	indexTemplate := template.Must(template.ParseFiles(indexFilepath))
	fmt.Printf("template : %v\n", indexTemplate)
	if err = indexTemplate.Execute(w, r); err != nil {
		fmt.Printf("Error occurs when serving the index page: %v\n", err.Error)
	}
}

func main() {
	// Serve all the files of the frontend.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Serve the request for the build list.
	http.HandleFunc("/buildList/", serveBuildListHandler)
	// Serve the request for the entity scores.
	http.HandleFunc("/entityScores/", serveEntityScoresHandler)
	// Serve the indexing page.
	http.HandleFunc("/", serveIndexHandler)
	// Serve all the other requests.
	http.ListenAndServe(":8000", nil)
}
