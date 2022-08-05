package main

import (
	"errors"
	"net/http"
)

func GetStrictCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value := cookie.Value
	if value == "" {
		return "", errors.New("value must be not empty")
	}

	return value, nil
}
