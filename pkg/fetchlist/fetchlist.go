package fetchlist

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Fetch fetches an url and returns all resulting lines
func Fetch(source string) ([]string, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, errors.New("error: " + res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(body), "\n"), nil
}
