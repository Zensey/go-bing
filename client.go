package bing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Token  string
	Client http.Client
}

//NewClient creates a new bing client istance
func NewClient(token string) *Client {
	return &Client{
		Token:  token,
		Client: *http.DefaultClient,
	}
}

//Simple Bing Search function
func (c *Client) Search(query *Query) (*BingAnswer, error) {
	//Build the request
	req, err := query.buildRequest()
	if err != nil {
		return nil, err
	}
	//set Header
	req.Header.Add("Ocp-Apim-Subscription-Key", c.Token)

	//send Request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status, err)

	//Retrieve request Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = checkBingErrors(body)
	if err != nil {
		return nil, err
	}
	//Parse Body json
	ans := new(BingAnswer)
	err = json.Unmarshal(body, &ans)
	return ans, err
}

func checkBingErrors(body []byte) error {
	berr := new(BingError)
	err := json.Unmarshal(body, &berr)
	if err != nil {
		return err
	}
	if berr.StatusCode == 401 {
		return fmt.Errorf(berr.Message)
	}
	return nil
}
