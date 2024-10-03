package cryptoapi

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// GetBitcoinPrice получает курс биткоина
func GetCryptoPrice(currency string) (float64, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%v&vs_currencies=usd", currency))

	if err != nil {
		return 0, err
	}

	var result map[string]map[string]float64
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, err
	}

	return result[currency]["usd"], nil
}
