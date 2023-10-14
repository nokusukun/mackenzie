package main

import (
	"encoding/json"
	"fmt"
	"github.com/nokusukun/mackenzie"
	"net/http"
	"time"
)

func HTTPGetJson(url string) (map[string]any, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	r := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	url := "https://dummyjson.com/products/1"
	CGetJson, err := mackenzie.Create[map[string]any](HTTPGetJson, mackenzie.Config{Lifetime: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	start := time.Now()
	data, _ := CGetJson.Get(url)
	elapsed := time.Since(start)
	fmt.Println("Data", data)
	fmt.Println("First call took", elapsed.String())

	start = time.Now()
	data, _ = CGetJson.Get(url)
	elapsed = time.Since(start)
	fmt.Println("Data", data)
	fmt.Println("Second call took", elapsed.String())

}
