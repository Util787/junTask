package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

// func to make concurent api requests with timeout
func requestUserAdditionalInfo(c *gin.Context, name string) (age int, gender string, nationality string) {

	// todo: test if they are reachable from goroutines in errgroup
	var ageResp agifyResponse
	var genderResp genderizeResponse
	var natResp nationalizeResponse

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errGr, ctx := errgroup.WithContext(timeOutCtx)

	errGr.Go(func() error {
		resp, err := requestUserAge(ctx, name)
		if err != nil {
			return err
		}
		ageResp = resp
		return nil
	})

	errGr.Go(func() error {
		resp, err := requestUserGender(ctx, name)
		if err != nil {
			return err
		}
		genderResp = resp
		return nil
	})

	errGr.Go(func() error {
		resp, err := requestUserNationality(ctx, name)
		if err != nil {
			return err
		}
		natResp = resp
		return nil
	})

	if err := errGr.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "requests time out or unreachable"})
		return
	}

	age = ageResp.Age
	gender = genderResp.Gender
	nationality = ""
	if len(natResp.Country) > 0 {
		nationality = natResp.Country[0].Country_id
	}
	return age, gender, nationality
}

func requestUserAge(ctx context.Context, name string) (agifyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", agifyURL, name), nil)
	if err != nil {
		return agifyResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Unreachable url ", agifyURL)
		return agifyResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp agifyResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		log.Println("Couldnt parse from ", agifyURL)
		return agifyResponse{}, err
	}
	return parsedResp, nil
}

func requestUserGender(ctx context.Context, name string) (genderizeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", genderizeURL, name), nil)
	if err != nil {
		return genderizeResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Unreachable url ", genderizeURL)
		return genderizeResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp genderizeResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		log.Println("Couldnt parse from ", genderizeURL)
		return genderizeResponse{}, err
	}
	return parsedResp, nil
}

func requestUserNationality(ctx context.Context, name string) (nationalizeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", nationalizeURL, name), nil)
	if err != nil {
		return nationalizeResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Unreachable url ", nationalizeURL)
		return nationalizeResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp nationalizeResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		log.Println("Couldnt parse from ", nationalizeURL)
		return nationalizeResponse{}, err
	}
	return parsedResp, nil
}
