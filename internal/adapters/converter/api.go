package converter

import (
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const templateURL = "https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%.2f"

// ExchangeRatesAPI is an example of currency conversion API
// This implementation should not be used in real production! Please, read about decimal values, currency conversions!
// https://apilayer.com/marketplace/exchangerates_data-api#documentation-tab
type ExchangeRatesAPI struct {
	apiKey      string
	templateURL string
	timeout     time.Duration
	client      *http.Client
}

func NewExchangeRatesAPI(apiKey string, timeout time.Duration) *ExchangeRatesAPI {
	client := http.Client{Timeout: timeout}
	a := ExchangeRatesAPI{
		apiKey:      apiKey,
		templateURL: templateURL,
		timeout:     timeout,
		client:      &client,
	}
	return &a
}

func (a *ExchangeRatesAPI) ConvertFromRUBToCurrency(amount float32, currency string) (float32, error) {
	url := a.createURL(amount, currency)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("apikey", a.apiKey)
	res, err := a.client.Do(req)
	if err != nil {
		return 0, err
	}
	var result apiResult
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	err = jsoniter.Unmarshal(bs, &result)
	if err != nil {
		return 0, err
	}
	return result.Result, nil
}

func (a *ExchangeRatesAPI) createURL(amount float32, currency string) string {
	return fmt.Sprintf(a.templateURL, currency, amount)
}
