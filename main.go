package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	//

	// buff := make([]byte, 512)
	// _, err = file.Read(buff)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// filetype := http.DetectContentType(buff)
	// if filetype != "application/json" {
	// 	http.Error(w, filetype+"The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
	// 	return
	// }

	// _, errr := file.Seek(0, io.SeekStart)
	// if errr != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Upload successful")

	http.Redirect(w, r, "/upload/complete", 302)

}


// SHOW MESSAGE THAT UPLOAD IS COMPLETED WITH AN OPTION TO SEARCH
func uploadComplete(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "uploadComplete.html")
}

// SHOW AIRLINE FORM
func searchAirline(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "searchAirline.html")
}

func notSuccededShowTable(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "notSuccededShowTable.html")
}

func ShowAirportResult(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("showUserTable.html")
		t.Execute(w, nil)

	} else {

		code := r.FormValue("id")
		fmt.Println(r.FormValue("id"))
		
		var alArpts AllAirports
		file, err := os.OpenFile("./uploads/airline.json", os.O_RDONLY, 0666)
		checkError(err)

		b, err := ioutil.ReadAll(file)
		checkError(err)

		json.Unmarshal(b, &alArpts.Airports)

		var allID []string
		for _, usr := range alArpts.Airports {
			allID = append(allID, usr.Airport.Code)
		}

		for _, usr := range alArpts.Airports {
			if IsValueInSlice(allID, code) != true {
				http.Redirect(w, r, "/notsuccessshowtable", 302)
				return
			}
			if usr.Airport.Code != code {
				continue
			} else {
				t, err := template.ParseFiles("showUserTable.html")
				checkError(err)
				t.Execute(w, usr)
			}
		}
	}
}


//MODEL FOR AIRPORTS
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func IsValueInSlice(slice []string, value string) (result bool) {
	for _, n := range slice {
		if n == value {
			return true
		}
	}
	return false
}

type Airports struct {
	Airport		*Airport	`json:"Airport"`
	Time 		*Time 		`json:"Time"`
	Statistics	*Statistics	`json:"Statistics"`
}

type Airport struct {
	Code	string	`json:"Code"`
	Name	string	`json:"Name"`
}

type Time struct{
	Label		string	`json:"Label"`
	Month		int		`json:"Month"`
	MonthName	string	`json:"Month Name"`
	Year		int		`json:"Year"`
}

type Statistics struct {
	Delays         	*Delays         `json:"# of Delays"`
	Flights        	*Flights        `json:"Flights"`
	MinutesDelayed	*MinutesDelayed	`json:"Minutes Delayed"`
}

type Delays struct {
	Carrier                int `json:"Minutes Delayed"`
	LateAircraft           int `json:"Late Aircraft"`
	NationalAviationSystem int `json:"National Aviation System "`
	Security               int `json:"Security"`
	Weather                int `json:"Weather"`
}

type Flights struct {
	Cancelled int `json:"Cancelled"`
	Delayed   int `json:"Delayed"`
	Diverted  int `json:"Diverted"`
	OnTime    int `json:"On Time"`
	Total     int `json:"Total"`
}

type MinutesDelayed struct {
	Carrier                int `json:"Carrier"`
	LateAircraft           int `json:"Late Aircraft"`
	NationalAviationSystem int `json:"National Aviation System"`
	Security               int `json:"Security"`
	Total                  int `json:"Total"`
	Weather                int `json:"Weather"`
}

type AllAirports struct {
	Airports []*Airports
}

//Show All Airports
func ShowAllAirports() (au *AllAirports) {
	file, err := os.OpenFile("./uploads/airline.json", os.O_RDWR|os.O_APPEND, 0666)
	checkError(err)
	b, err := ioutil.ReadAll(file)
	var alArpts AllAirports
	json.Unmarshal(b, &alArpts.Airports)
	checkError(err)
	return &alArpts
}

func main() {

	http.HandleFunc("/", indexHandler)

	//Uploading and error handling pages
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/upload/complete", uploadComplete)

	http.HandleFunc("/searchAirline", searchAirline)

	http.HandleFunc("/airlinetable", ShowAirportResult)
	http.HandleFunc("/notsuccessshowtable", notSuccededShowTable)


	fmt.Println("Listening...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
