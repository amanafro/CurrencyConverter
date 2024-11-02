package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

var (
	BaseCurrency     string
	CovertedCurrency string
	Amount           float64
	AmountStr        string
	CovertedAmount   float64
	confirm          bool
)

func main() {

	// godotenv.Load(".env")
	appID, exists := os.LookupEnv("API_KEY")
	if !exists {
		log.Fatalln("API key not found in environment variables. Please set the API_KEY environment variable.")
	}

	fmt.Println("API Key:", appID)
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", appID)

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

	baseRate := exchangeRates.Rates[BaseCurrency]
	convertedRate := exchangeRates.Rates[CovertedCurrency]
	CovertedAmount := Amount * (convertedRate / baseRate)

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
				Title("Enter the amount").
				Value(&AmountStr). // Store in amountStr instead of amount
				Validate(func(str string) error {
					val, err := strconv.ParseFloat(str, 64)
					if err != nil || val <= 0 {
						return errors.New("please enter a valid amount")
					}
					Amount = val // Assign to amount if parsing is successful
					return nil
				}),

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
