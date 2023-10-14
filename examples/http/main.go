package main

import (
	"encoding/json"
	"fmt"
	"github.com/nokusukun/mackenzie"
	"net/http"
	"time"
)

type Response struct {
	Id                 int      `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Price              int      `json:"price"`
	DiscountPercentage float64  `json:"discountPercentage"`
	Rating             float64  `json:"rating"`
	Stock              int      `json:"stock"`
	Brand              string   `json:"brand"`
	Category           string   `json:"category"`
	Thumbnail          string   `json:"thumbnail"`
	Images             []string `json:"images"`
}

func HTTPGetJson[T any](url string) (T, error) {
	r := new(T)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return *r, err
	}
	resp, err := client.Do(req)
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return *r, err
	}
	return *r, nil
}

func main() {
	url := "https://dummyjson.com/products/1"
	CGetJson, err := mackenzie.Create[Response](HTTPGetJson[Response], mackenzie.Config{Lifetime: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	start := time.Now()
	data, _ := CGetJson.Get(url)
	elapsed := time.Since(start)
	fmt.Println("Data", data.Id)
	fmt.Println("First call took", elapsed.String())

	start = time.Now()
	data, _ = CGetJson.Get(url)
	elapsed = time.Since(start)
	fmt.Println("Data", data.Id)
	fmt.Println("Second call took", elapsed.String())

}
