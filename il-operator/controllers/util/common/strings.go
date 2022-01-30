package common

import (
	"fmt"
	"strings"
)

// ContainsString is a helper functions to check and remove string from a slice of strings.
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func Stringify(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Find(s []string, e string) *string {
	for _, a := range s {
		if a == e {
			return &e
		}
	}
	return nil
}

func ExtractNameFromPath(path string) string {
	tokens := strings.Split(path, "/")
	return tokens[len(tokens)-1]
}

func ToStringArray(stringArr []string) string {
	arr := ""
	for i, s := range stringArr {
		if i != 0 {
			s += ","
		}
		arr += s
	}
	return arr
}
