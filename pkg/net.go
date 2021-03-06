package pkg

import (
	"fmt"
	"io"
	"net/http"
)

func Get200ResBody(c *http.Client, r *http.Request) ([]byte, error) {
	res, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected response status 200, got %d", res.StatusCode)
	}

	return io.ReadAll(res.Body)
}

func GetResBody(c *http.Client, r *http.Request) ([]byte, int, error) {
	res, err := c.Do(r)
	if err != nil {
		return nil, -1, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	return data, res.StatusCode, err
}
