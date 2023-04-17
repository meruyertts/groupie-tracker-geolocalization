package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Id             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      string              `json:"locations"`
	ConcertDates   string              `json:"concertDates"`
	Relations      string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ErrorBody struct {
	Status  int
	Message string
}

var All *template.Template

func init() {
	var err error
	All, err = template.ParseFiles("ui/htmlfiles/error.html", "ui/htmlfiles/index.html", "ui/htmlfiles/more-info.html")
	if err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHeader(w, http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		errorHeader(w, http.StatusMethodNotAllowed)
		return
	}
	var result []Response
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err.Error())
		return
	}

	err = All.ExecuteTemplate(w, "index.html", result)
	if err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	var result Response
	id := strings.TrimPrefix(r.URL.Path, "/artist/")
	_, err := strconv.Atoi(id)
	if err != nil {
		errorHeader(w, http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		errorHeader(w, http.StatusMethodNotAllowed)
		return
	}
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
	if result.Id == 0 {
		errorHeader(w, http.StatusNotFound)
		return
	}
	relData, err := http.Get(result.Relations)
	if err != nil {

		errorHeader(w, http.StatusInternalServerError)
		return
	}
	defer relData.Body.Close()

	if err = json.NewDecoder(relData.Body).Decode(&result); err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
	err = All.ExecuteTemplate(w, "more-info.html", result)
	if err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
}

func errorHeader(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	errH := setError(status)
	err := All.ExecuteTemplate(w, "error.html", errH)
	if err != nil {
		errorHeader(w, http.StatusInternalServerError)
		return
	}
}

func setError(status int) *ErrorBody {
	return &ErrorBody{
		Status:  status,
		Message: http.StatusText(status),
	}
}
