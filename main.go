package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

// TaxBracket represents a single tax bracket
type TaxBracket struct {
	Rate float64  `json:"rate"`
	UpTo *float64 `json:"up_to"` // nil means no upper limit (last bracket)
}

// TaxConfig holds multiple tax brackets
type TaxConfig struct {
	Brackets []TaxBracket `json:"brackets"`
}

// LoadConfig reads JSON config from a file
func LoadConfig(filename string) (*TaxConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config TaxConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Sort brackets by income limit
	sort.Slice(config.Brackets, func(i, j int) bool {
		if config.Brackets[i].UpTo == nil {
			return false
		}
		if config.Brackets[j].UpTo == nil {
			return true
		}
		return *config.Brackets[i].UpTo < *config.Brackets[j].UpTo
	})

	return &config, nil
}

// CalculateTax computes tax for given income
func CalculateTax(income float64, brackets []TaxBracket) (float64, map[float64]float64) {
	totalTax := 0.0
	taxDetails := make(map[float64]float64)
	previousLimit := 0.0

	for _, bracket := range brackets {
		taxableAmount := income - previousLimit
		if bracket.UpTo != nil && taxableAmount > (*bracket.UpTo-previousLimit) {
			taxableAmount = *bracket.UpTo - previousLimit
		}

		if taxableAmount > 0 {
			taxForBracket := taxableAmount * (bracket.Rate / 100)
			totalTax += taxForBracket
			taxDetails[bracket.Rate] = taxForBracket
		}

		previousLimit = *bracket.UpTo
		if bracket.UpTo == nil || income <= previousLimit {
			break
		}
	}

	return totalTax, taxDetails
}

func main() {
	var cfile string
	flag.StringVar(&cfile, "config", "config.json", "tax configuration file")
	flag.Parse()

	// Load tax brackets from JSON file
	config, err := LoadConfig(cfile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Example income
	var income float64
	fmt.Print("Enter your taxable income: ")
	fmt.Scan(&income)

	// Calculate tax
	totalTax, taxBreakdown := CalculateTax(income, config.Brackets)

	// Print results
	fmt.Println("\nTax Breakdown:")
	for rate, tax := range taxBreakdown {
		fmt.Printf("- %.0f%%: $%.2f\n", rate, tax)
	}

	fmt.Printf("\nTotal Tax Owed: $%.2f\n", totalTax)
	fmt.Printf("Effective Tax Rate: %.2f%%\n", (totalTax/income)*100)
}
