package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type registerResponse struct {
	PlayerId int
}

type collisionResponse struct {
	Role string
}

func getResponse(url string, response interface{}) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to communicate to url '%s', received error: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	return dec.Decode(response)
}
