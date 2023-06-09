package utils

func IsStringInSlice(stringToFind string, slice []string) bool {
	for _, stringInSlice := range slice {
		if stringInSlice == stringToFind {
			return true
		}
	}
	return false
}
