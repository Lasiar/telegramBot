package system

func DeleteByValue(value int64, array []int64) []int64 {
	var arrayReturn []int64
	for _, a := range array {
		if value != a {
			arrayReturn = append(arrayReturn, a)
		}
	}
	return arrayReturn
}