package converter

// Converter TODO: move to controllers/usecase
type Converter interface {
	ConvertFromRUBToCurrency(currency string) (float32, error)
}
