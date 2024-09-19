package main

import (
	"errors"
	"github.com/charmbracelet/huh"
	"log"
)

var (
	BaseCurrency     string
	CovertedCurrency string
	Amount           string
	instructions     string
	discount         bool
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewSelect[string]().
				Title("Choose your burger").
				Options(
					huh.NewOption("Fr CHF", "CHF"),
					huh.NewOption("€ EUR", "EUR"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
				).
				Value(&BaseCurrency), // store the chosen option in the "burger" variable

			// Let the user select multiple toppings.
			huh.NewSelect[string]().
				Title("Toppings").
				Options(
					huh.NewOption("Fr CHF", "CHF"),
					huh.NewOption("€ EUR", "EUR"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
				).
				Value(&CovertedCurrency),
			huh.NewInput().
				Title("How much would you want to").
				Value(&Amount).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if str == "Frank" {
						return errors.New("Sorry, we don’t serve customers named Frank.")
					}
					return nil
				}),
		),

		// Gather some final details about the order.
		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to confirm?").
				Value(&discount),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
