package converter

import (
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

// const templateURL = "https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%s"

// exchangeRatesAPI is an example of currency conversion API
// This implementation should not be used in real production! Please, read about decimal values, currency conversions!
// https://apilayer.com/marketplace/exchangerates_data-api#documentation-tab
type exchangeRatesAPI struct {
	apiKey      string
	apiURL      string
	templateURL string
	timeout     time.Duration
	client      *http.Client
}

func NewExchangeRatesAPI(apiKey, apiURL string, timeout time.Duration) *exchangeRatesAPI {
	client := http.Client{Timeout: timeout}
	templateURL := apiURL + "/convert?to=%s&from=RUB&amount=%s"
	a := exchangeRatesAPI{
		apiKey:      apiKey,
		apiURL:      apiURL,
		templateURL: templateURL,
		timeout:     timeout,
		client:      &client,
	}
	return &a
}

func (a *exchangeRatesAPI) ConvertFromRUBToCurrency(amount decimal.Decimal, currency string) (decimal.Decimal, error) {
	url := a.createURL(amount.String(), currency)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return decimal.Zero, err
	}
	req.Header.Set("apikey", a.apiKey)
	res, err := a.client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	var result apiResult
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return decimal.Zero, err
	}
	err = jsoniter.Unmarshal(bs, &result)
	if err != nil {
		return decimal.Zero, err
	}
	// TODO
	return result.Result, nil
}

func (a *exchangeRatesAPI) createURL(amount string, currency string) string {
	return fmt.Sprintf(a.templateURL, currency, amount)
}
