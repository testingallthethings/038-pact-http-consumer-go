package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type JsonError struct {
	Code string `json:"code" pact:"example=1234"`
	Msg  string `json:"msg" pact:"example=No book with ISBN 123456789"`
}

type Book struct {
	ISBN          string `json:"isbn" pact:"example=987654321"`
	Title         string `json:"title" pact:"example=Testing All The Things"`
	Image         string `json:"image" pact:"example=testing.jpg"`
	Genre         string `json:"genre" pact:"example=Computers"`
	YearPublished int    `json:"year_published" pact:"example=2021"`
}

func NewClient(host string) Client {
	return Client{host}
}

type Client struct {
	host string
}

func (c Client) GetBook(isbn string) (Book, error) {
	client := http.Client{}

	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/book/%s", c.host, isbn),
		nil,
		)
	req.Header.Add("Accept", "application/json")

	resp, _ := client.Do(req)

	if resp.StatusCode == http.StatusNotFound {
		je := JsonError{}
		json.NewDecoder(resp.Body).Decode(&je)

		return Book{}, errors.New(fmt.Sprintf("%s - %s", je.Code, je.Msg))
	}

	b := Book{}
	json.NewDecoder(resp.Body).Decode(&b)

	return b, nil

}
