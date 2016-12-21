package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fjukstad/met"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/snow", SnowHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err)
		return
	}
}

type Observation struct {
	Thickness float64
	Date      time.Time
}

func SnowHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	year := v.Get("year")
	if year == "" {
		http.Error(w, "no year specified", 400)
		return
	}
	f := met.Filter{
		Sources:       []string{"SN90450"},
		ReferenceTime: year + "-01-01T00:00:00.000Z/" + year + "-12-30T00:00:00.000Z",
		Elements:      []string{"surface_snow_thickness"},
	}

	data, err := met.GetObservations(f)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	observations := []Observation{}

	for _, d := range data {
		thickness := d.Observations[0].Value
		date := d.ReferenceTime
		observations = append(observations, Observation{thickness, date})

	}

	b, err := json.Marshal(observations)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("index.html")
	if err != nil {
		fmt.Println("Could not read index.html")
		http.Error(w, "not found", 404)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Could not read index.html")
		http.Error(w, "index.html not found", 404)
		return
	}

	w.Write(b)
}
