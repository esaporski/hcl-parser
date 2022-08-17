package main

// Tells if array contains string
func Contains(arr []string, str string) bool {
	for _, item := range arr {
		if str == item {
			return true
		}
	}
	return false
}
