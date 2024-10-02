package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
)

// WARNING: this code is shit

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

var (
	BaseCurrency     string
	CovertedCurrency string
	Amount           string
	CovertedAmount   float64
	confirm          bool
)

func main() {

	// godotenv.Load(".env")
	appID, exists := os.LookupEnv("API_KEY")
	if !exists {
		log.Fatalln("API key not found in environment variables. Please set the API_KEY environment variable.")
	}

	// FIX: the conversion vaule is not showing and i'm getting NaN as a response

	fmt.Println("API Key:", appID)
	url := fmt.Sprint("https://openexchangerates.org/api/latest.json?app_id=%s", appID)

	// Make the HTTP request
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Decode of the json

	var exchangeRates ExchangeRates
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: i had to parse the string into float64, because for some reason the value() function only accepts string and i was lazy to not ckeck or try to find a solution

	amountFloat, err := strconv.ParseFloat(Amount, 64)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Converting logic

	baseRate := exchangeRates.Rates[BaseCurrency]
	convertedRate := exchangeRates.Rates[CovertedCurrency]
	CovertedAmount := amountFloat * (convertedRate / baseRate)

	fmt.Printf("Converted amout: %.2f %s\n", CovertedAmount, CovertedCurrency)

	// TODO: add component to display the converted amount and currency
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your base currency").
				Options(
					huh.NewOption("Fr CHF", "CHF"),
					huh.NewOption("€ EUR", "EUR"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
				).
				Value(&BaseCurrency),

			huh.NewSelect[string]().
				Title("What do you want it to convert into?").
				Options(
					huh.NewOption("Fr CHF", "CHF"),
					huh.NewOption("€ EUR", "EUR"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
				).
				Value(&CovertedCurrency),
			huh.NewInput().
				Title("How much would you want to").
				Value(&Amount),

			huh.NewConfirm().
				Title("Would you like to confirm?").
				Value(&confirm),
		),
	)

	// NOTE: Execute the form
	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
