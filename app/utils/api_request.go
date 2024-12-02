package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func RequestAPI(method string, url string, data interface{}, out interface{}) error {
	// Marshal the input data to JSON
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Initialize HTTP client and send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close() // Ensure body is closed after reading

	// Check for non-2xx status codes
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return errors.New("request failed with status " + res.Status + ": " + string(bodyBytes))
	}

	// Read and unmarshal the response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	logrus.Printf("Response: %s", string(bodyBytes))

	// Unmarshal response JSON into output struct
	err = json.Unmarshal(bodyBytes, out)
	if err != nil {
		return err
	}

	return nil
}
