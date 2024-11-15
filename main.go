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
	"github.com/charmbracelet/lipgloss"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type CurrencyConverter struct {
	BaseCurrency     string
	CovertedCurrency string
	Amount           float64
	AmountStr        string
	CovertedAmount   float64
	confirm          bool
	rates            ExchangeRates
}

func (cc *CurrencyConverter) convertCurrency() error {
	baseRate, ok := cc.rates.Rates[cc.BaseCurrency]
	if !ok {
		return fmt.Errorf("base currency %s not found", cc.BaseCurrency)
	}
	convertedRate, ok := cc.rates.Rates[cc.CovertedCurrency]
	if !ok {
		return fmt.Errorf("base currency %s not found", cc.CovertedCurrency)
	}
	cc.CovertedAmount = cc.Amount * (convertedRate / baseRate)
	return nil
}

func fetchExchangeRate() (ExchangeRates, error) {
	appID, exists := os.LookupEnv("API_KEY")
	if !exists {
		return ExchangeRates{}, errors.New("API key not found in .env")
	}

	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", appID)

	// Make the HTTP request
	response, err := http.Get(url)
	if err != nil {
		return ExchangeRates{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ExchangeRates{}, fmt.Errorf("API request failed with status: %d", response.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ExchangeRates{}, err
	}

	var exchangeRates ExchangeRates

	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		return ExchangeRates{}, err
	}

	return exchangeRates, nil
}

func main() {

	converter := &CurrencyConverter{}
	rates, err := fetchExchangeRate()

	if err != nil {
		log.Fatal(err)
	}

	converter.rates = rates

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
				Value(&converter.BaseCurrency),

			huh.NewSelect[string]().
				Title("What do you want it to convert into?").
				Options(
					huh.NewOption("Fr CHF", "CHF"),
					huh.NewOption("€ EUR", "EUR"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
				).
				Value(&converter.CovertedCurrency),

			huh.NewInput().
				Title("Enter the amount").
				Value(&converter.AmountStr). // Store in amountStr instead of amount
				Validate(func(str string) error {
					val, err := strconv.ParseFloat(str, 64)
					if err != nil || val <= 0 {
						return errors.New("please enter a valid amount")
					}
					converter.Amount = val // Assign to amount if parsing is successful
					return nil
				}),

			huh.NewConfirm().
				Title("Would you like to confirm?").
				Value(&converter.confirm),
		),
	)

	// NOTE: Execute the form
	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// FIX: The display for the converted amount looks out of place

	if converter.confirm {
		if err := converter.convertCurrency(); err != nil {
			log.Fatal(err)
		}
		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingLeft(4).
			Width(30)
		s := fmt.Sprintf("Converted amount:%s %0.2f", converter.CovertedCurrency, converter.CovertedAmount)
		fmt.Println(style.Render(s))
	}
}
