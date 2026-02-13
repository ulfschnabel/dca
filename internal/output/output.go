package output

import (
	"encoding/json"
	"fmt"
)

// Response is the standard JSON response format
type Response struct {
	OK    bool        `json:"ok"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// Success creates a successful response
func Success(data interface{}) *Response {
	return &Response{
		OK:   true,
		Data: data,
	}
}

// Error creates an error response
func Error(err error) *Response {
	return &Response{
		OK:    false,
		Error: err.Error(),
	}
}

// Print outputs the response as JSON
func Print(resp *Response, pretty bool) error {
	var data []byte
	var err error

	if pretty {
		data, err = json.MarshalIndent(resp, "", "  ")
	} else {
		data, err = json.Marshal(resp)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// PrintSuccess prints a successful response
func PrintSuccess(data interface{}, pretty bool) error {
	return Print(Success(data), pretty)
}

// PrintError prints an error response
func PrintError(err error, pretty bool) error {
	return Print(Error(err), pretty)
}
