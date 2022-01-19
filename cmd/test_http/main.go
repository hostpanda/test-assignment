package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	values := map[string]string{"event": "test", "occupation": "gardener"}

	data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8085/addHTTP/", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
}
