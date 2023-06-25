package entity

import "net/http"

var (
	CLI = &http.Client{}
)

const (
	API_HOST  = "API_HOST"
	API_KEY   = "API_KEY"
	MAX_TOKEN = "MAX_TOKEN"
)
