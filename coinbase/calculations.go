package coinbase

// Calculate average of an array of values
func Average(values []float64) (result float64) {
	for _, value := range values {
		result += value
	}
	result = result / float64(len(values))

	return
}

// Calculate percentage with respect to the Current Price
func GetPercentageDifference(price, value float64) float64 {
	return ((price - value) * 100 / value)
}
