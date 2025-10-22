package ingest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Base  string
	Token string
	HTTP  *http.Client
}

func New(base, token string) *Client { return &Client{Base: base, Token: token, HTTP: &http.Client{}} }

type Page struct {
	Items any    `json:"items"`
	Next  string `json:"next_page"`
}

func (c *Client) Fetch(next string) (Page, error) {
	u, _ := url.Parse(c.Base)
	if next != "" {
		q := u.Query()
		q.Set("next_page", next)
		u.RawQuery = q.Encode()
	}
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Page{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return Page{}, errors.New(fmt.Sprintf("status %d: %s", resp.StatusCode, string(b)))
	}
	var p Page
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&p); err != nil {
		return Page{}, err
	}
	return p, nil
}
