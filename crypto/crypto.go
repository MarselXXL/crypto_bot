package crypto

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

// GetBitcoinPrice получает курс биткоина
func GetBitcoinPrice() (float64, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")

	if err != nil {
		return 0, err
	}

	var result map[string]map[string]float64
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, err
	}

	return result["bitcoin"]["usd"], nil
}
