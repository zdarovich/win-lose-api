package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zdarovich/win-lose-api/handlers"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL   *url.URL
	httpClient *http.Client
}

func (c *Client) call(body *handlers.TrxRequest) (interface{}, error) {
	rel := &url.URL{Path: "/your_url"}
	u := c.BaseURL.ResolveReference(rel)
	
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Source-Type", "game")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		var trx handlers.TrxResponse
		err = json.NewDecoder(resp.Body).Decode(&trx)
		return &trx, err
	} else {
		type ErrorResp struct {
			Error string  `json:"error"`
		}
		var errResp ErrorResp
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		return &errResp, err
	}
	
}


func main() {
	timeout := time.After(60 * time.Second)
	tick := time.Tick(1 * time.Second)
	c := &Client {BaseURL: &url.URL{
			Scheme:     "http",
			Host:       "127.0.0.1:8081",
		},
		httpClient: http.DefaultClient,
	}
	max := 10.5
	min := 0.0


	for {
		select {
		case <-timeout:
			log.Fatal("timed out")
			return
		case <-tick:
			r := min + rand.Float64() * (max - min)
			var state string
			if rand.Float32() < 0.5 {
				state = "lost"
			} else {
				state = "win"
			}
			u := uuid.New().String()
			resp, err := c.call(&handlers.TrxRequest{
				State:         state,
				Amount:        fmt.Sprintf("%f", r),
				TransactionId: u,
			})
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Printf("%+v \n", resp)

		}
	}
}