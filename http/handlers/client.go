package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type ClientHandlers struct{}

const authUrl = "http://localhost:8080/auth/check"

func NewClientHandlers() *ClientHandlers {
	return &ClientHandlers{}
}

func (ch *ClientHandlers) Ping(rw http.ResponseWriter, r *http.Request) {
	log.Println("Ping Request")

	header := r.Header.Get("Authorization")
	token, err := parseToken(header)
	if err != nil {
		http.Error(rw, "Bad Authorization", http.StatusUnauthorized)
		return
	}

	payload, err := json.Marshal(map[string]string{
		"token": token,
	})

	body := bytes.NewBuffer(payload)

	res, err := http.Post(authUrl, "application/json", body)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	rw.Header().Add("Content-Type", "application/json")

	_, err = io.Copy(rw, res.Body)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func parseToken(str string) (string, error) {
	arr := strings.Split(str, " ")

	if len(arr) != 2 {
		return "", fmt.Errorf("bad token")
	}

	if arr[0] != "Bearer" {
		return "", fmt.Errorf("bad token")
	}

	return arr[1], nil
}
