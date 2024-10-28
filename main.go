package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
	"os"
	"text/template"
	"strings"
)

import "github.com/stianeikeland/go-rpio/v4"

func main() {
    	
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) { start(w) })
	http.HandleFunc("/on", func(w http.ResponseWriter, r *http.Request) { on() })
	http.HandleFunc("/off", func(w http.ResponseWriter, r *http.Request) { off() })
	http.HandleFunc("/trivia", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "%s° F", trivia()) })
	http.HandleFunc("/movie/", GetMovie)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/movie/")
	const tplText = `
	<video width="320" height="240" controls>
	  <source src="/movies/{{.}}" type="video/mp4">
	Your browser does not support the video tag.
	</video>
		   	`


	tpl, err := template.New("nameList").Parse(tplText)
    	if err != nil {
        	fmt.Println("Error parsing template:", err)
        	return
    	}

    	// Execute the template with the array of names
    	err = tpl.Execute(w, name)
    	if err != nil {
        	fmt.Println("Error executing template:", err)
        	return 
	}
	
}

func start(w http.ResponseWriter) {
	f, err := os.Open("/home/rschofield/homepage/static/movies")
	
    	if err != nil {
        	fmt.Println(err)
        	return
    	}
    	files, err := f.Readdirnames(0)
    	if err != nil {
        	fmt.Println(err)
        	return
    	}

	//content, err := os.ReadFile("/home/rschofield/homepage/static/video-player.html")

	//if err != nil {
	//	return "wtf"
	//}

	const tplText = `
		<ul>
			{{range .}}
				<li hx-get="/movie/{{.}}" hx-target="#video" hx-trigger="click">{{.}}</li>
			{{end}}
		</ul>
   	`
	tpl, err := template.New("nameList").Parse(tplText)
    	if err != nil {
        	fmt.Println("Error parsing template:", err)
        	return
    	}

    	// Execute the template with the array of names
    	err = tpl.Execute(w, files)
    	if err != nil {
        	fmt.Println("Error executing template:", err)
        	return 
	}


}

func on() {
	err := rpio.Open()

	if err == nil {
		pin := rpio.Pin(17)
		pin.Output()
		pin.High()
		rpio.Close()
	}
}

func off() {
	err := rpio.Open()

	if err == nil {
		pin := rpio.Pin(17)
		pin.Output()
		pin.Low()
		rpio.Close()
	}
}

type Forecast struct {
	Properties struct {
		Periods []struct {
			IsDaytime   bool   `json:"isDaytime"`
			Temperature int    `json:"temperature"`
			ProbabilityOfPrecipitation struct {
				UnitCode string `json:"unitCode"`
				Value    int    `json:"value"`
			} `json:"probabilityOfPrecipitation"`
			RelativeHumidity struct {
				UnitCode string `json:"unitCode"`
				Value    int    `json:"value"`
			} `json:"relativeHumidity"`
			WindSpeed        string `json:"windSpeed"`
			WindDirection    string `json:"windDirection"`
			Icon             string `json:"icon"`
			ShortForecast    string `json:"shortForecast"`
			DetailedForecast string `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func trivia() string {
	resp, err := http.Get("https://api.weather.gov/gridpoints/TOP/98,47/forecast")

	if err != nil {
		return "wtf"
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "wtf"
	}

	var data Forecast
	err2 := json.Unmarshal(body, &data)
	if err2 != nil {
		return "wtf"
	}
	return strconv.Itoa(data.Properties.Periods[0].Temperature)
}
