package ean

import "regexp"

func IsValid(ean string) bool {
	ean8Regex := regexp.MustCompile("^[0-9]{8}$")
	upcARegex := regexp.MustCompile("^[0-9]{12}$")
	ean13Regex := regexp.MustCompile("^[0-9]{13}$")

	return ean8Regex.MatchString(ean) || ean13Regex.MatchString(ean) || upcARegex.MatchString(ean)
}
