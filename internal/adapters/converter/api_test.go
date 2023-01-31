package converter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestConvertFromRUBToCurrency_Success(t *testing.T) {
	var expectedErr error
	expectedResult, e := decimal.NewFromString("6.4")
	if e != nil {
		assert.Fail(t, fmt.Sprintf("expectedResult creating error is not nil; %v", e))
	}

	jsonString := `{
    "success": true,
    "query": {
        "from": "RUB",
        "to": "JPY",
        "amount": 2
    },
    "info": {
        "timestamp": 1519328414,
        "rate": 3.2
    },
    "historical": "",
    "date": "2018-02-22",
    "result": 6.4
}`
	amount := decimal.NewFromInt(2)
	currency := "JPY"
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(jsonString))
			},
		),
	)
	defer server.Close()
	converter := NewExchangeRatesAPI("", server.URL, time.Second)
	result, err := converter.ConvertFromRUBToCurrency(amount, currency)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedResult, result)
}
func TestConvertFromRUBToCurrency_JSONDecodeError(t *testing.T) {
	expectedResult := decimal.Zero
	amount := decimal.NewFromInt(2)
	currency := "JPY"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	converter := NewExchangeRatesAPI("", server.URL, time.Second)
	result, err := converter.ConvertFromRUBToCurrency(amount, currency)
	assert.ErrorContains(t, err, "readObjectStart: expect { or n")
	assert.Equal(t, expectedResult, result)
}
