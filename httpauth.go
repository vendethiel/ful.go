package main

/**
 * http-related helpers
 */

import (
	"encoding/base64"
	"strings"
)

type httpCode int

type BasicUser struct {
	Username, Password string
}

func parseAuthHeader(authHeader []string) (*BasicUser, httpCode) {
	// validate we have the header... (not present: unauthorized)
	if len(authHeader) == 0 {
		return nil, 401
	}
	// validate the header's shape... (bad shape: bad request)
	auth := strings.SplitN(authHeader[0], " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return nil, 400
	}

	// parse the header... validate the size (bad shape: bad request)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return nil, 400
	}
	return &BasicUser{pair[0], pair[1]}, 200
}
