package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type bTime struct {
	Updated    string
	UpdatedISO string
	Updateduk  string
}

type bEur struct {
	Code        string
	Rate        string
	Description string
	Rate_float  float64
}

type bCzk struct {
	Code        string
	Rate        string
	Description string
	Rate_float  float64
}

type bUsd struct {
	Code        string
	Rate        string
	Description string
	Rate_float  float64
}

type bBpi struct {
	Usd bUsd
	Eur bEur
	Czk bCzk
}

type bResponse struct {
	Time       bTime  `json:"time"`
	Disclaimer string `json:"disclaimer"`
	Rate       string `json:"rate"`
	Bpi        bBpi   `json:"bpi"`
	Date       time.Time
}

type result struct {
	CodeEur       string    `json:"codeEur"`
	RateEur       string    `json:"rateEur"`
	UpdatedEur    string    `json:"updatedEur"`
	ClientTimeEur time.Time `json:"clienttimeEur"`
	CodeCzk       string    `json:"codeCzk"`
	RateCzk       string    `json:"rateCzk"`
	UpdatedCzk    string    `json:"updatedCzk"`
	ClientTimeCzk time.Time `json:"clienttimeczk"`
}

func main() {
	// Start and run web server on port 8080
	http.HandleFunc("/", Handler)
	fmt.Println("Web server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}

// Request api.coindesk.com/ for price.
func getPrice(symbol string) (string, string, time.Time, string, error) {
	url := fmt.Sprintf("https://api.coindesk.com/v1/bpi/currentprice/%s.json", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "error", "error", time.Time{}, "error", err
	}

	var bRes bResponse

	err = json.NewDecoder(resp.Body).Decode(&bRes)
	if err != nil {
		return "error", "error", time.Time{}, "error", err
	}

	bRes.Date = time.Now()

	if symbol == "EUR" {

		return bRes.Bpi.Eur.Code, bRes.Time.Updated, bRes.Date, bRes.Bpi.Eur.Rate, err

	}

	return bRes.Bpi.Czk.Code, bRes.Time.Updated, bRes.Date, bRes.Bpi.Czk.Rate, err
}

func Handler(w http.ResponseWriter, r *http.Request) {

	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		fmt.Println(error)
	}

	defer conn.Close()

	// Some hardcoding
	eur := "EUR"
	czk := "CZK"

	codeeur, updatedeur, clienttimeeur, rateeur, _ := getPrice(eur)
	codeczk, updatedczk, clienttimeczk, rateczk, _ := getPrice(czk)

	result := result{
		CodeEur:       codeeur,
		RateEur:       rateeur,
		UpdatedEur:    updatedeur,
		ClientTimeEur: clienttimeeur,
		CodeCzk:       codeczk,
		RateCzk:       rateczk,
		UpdatedCzk:    updatedczk,
		ClientTimeCzk: clienttimeczk,
	}

	fmt.Println(result)

	// Marshal the result struct to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	jsonString := string(jsonData)
	fmt.Println("JSON Data:")
	fmt.Println(jsonString)

	io.WriteString(w, (jsonString))

}
