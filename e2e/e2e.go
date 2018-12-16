package e2e

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	BaseURL  = "http://127.0.0.1:8080"
	UserName = "heroku"
	Password = "plivo"
)

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func HttpRequest(reqType string, url string, headers map[string]string, body io.Reader) (res []byte, err error) {
	var req *http.Request
	client := &http.Client{}

	// Create request
	req, err = http.NewRequest(reqType, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Basic "+BasicAuth(UserName, Password))
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
