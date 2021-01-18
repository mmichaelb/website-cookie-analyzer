package util

// GetMedian returns the median of the given input data. The data has to be sort in order for the function to return a
// correct value.
func GetMedian(data []int) float32 {
	dataLength := len(data)
	if dataLength == 0 {
		return -1
	}
	middle := dataLength / 2
	// middle present
	if dataLength == 1 {
		return float32(data[0])
	} else if dataLength%2 == 1 {
		return float32(data[middle])
	} else {
		summed := data[middle] + data[middle-1]
		return float32(float64(summed) / 2)
	}
}
