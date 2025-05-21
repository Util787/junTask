package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	agifyURL       = "https://api.agify.io/"
	genderizeURL   = "https://api.genderize.io/"
	nationalizeURL = "https://api.nationalize.io/"
)

type agifyResponse struct {
	Age int `json:"age"`
}

type genderizeResponse struct {
	Gender string `json:"gender"`
}

type nationalizeResponse struct {
	Country []countryInfo `json:"country"`
}

type countryInfo struct {
	Country_id  string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

// func to make api requests
func requestUserAdditionalInfo(c *gin.Context, name string) {

	// the point of making a buffer in chans is to avoid getting goroutine running forever after context is done

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
}

func requestUserAge(ctx context.Context, name string) (agifyResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", agifyURL, name))
	if err != nil {
		return agifyResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp agifyResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return agifyResponse{}, err
	}
	return parsedResp, nil
}

func requestUserGender(ctx context.Context, name string) (genderizeResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", genderizeURL, name))
	if err != nil {
		return genderizeResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp genderizeResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return genderizeResponse{}, err
	}
	return parsedResp, nil
}

func requestUserNationality(ctx context.Context, name string) (nationalizeResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", nationalizeURL, name))
	if err != nil {
		return nationalizeResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp nationalizeResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return nationalizeResponse{}, err
	}
	return parsedResp, nil
}
