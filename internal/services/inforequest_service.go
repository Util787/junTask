package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const apiCallsTimeOut = 5 * time.Second

func NewInfoRequestService() InfoRequestService {
	return &infoRequestService{}
}

type infoRequestService struct{}

func (r *infoRequestService) RequestAdditionalInfo(name string) (int, string, string, error) {
	return requestUserAdditionalInfo(name)
}

const (
	agifyURL       = "https://api.agify.io/"
	genderizeURL   = "https://api.genderize.io/"
	nationalizeURL = "https://api.nationalize.io/"
)

type agifyResponse struct {
	Age   int    `json:"age"`
	Error string `json:"error"`
}

type genderizeResponse struct {
	Gender string `json:"gender"`
	Error  string `json:"error"`
}

type nationalizeResponse struct {
	Country []countryInfo `json:"country"`
	Error   string        `json:"error"`
}

type countryInfo struct {
	Country_id  string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

// func to make concurent api requests with timeout
func requestUserAdditionalInfo(name string) (age int, gender string, nationality string, err error) {

	var ageResp agifyResponse
	var genderResp genderizeResponse
	var natResp nationalizeResponse

	timeOutCtx, cancel := context.WithTimeout(context.Background(), apiCallsTimeOut)
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
		return 0, "", "", err
	}

	age = ageResp.Age
	gender = genderResp.Gender
	nationality = ""
	if len(natResp.Country) > 0 {
		nationality = natResp.Country[0].Country_id
	}
	return age, gender, nationality, nil
}

func requestUserAge(ctx context.Context, name string) (agifyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", agifyURL, name), nil)
	if err != nil {
		return agifyResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return agifyResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp agifyResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return agifyResponse{}, err
	}
	if parsedResp.Error != "" {
		return agifyResponse{}, errors.New(parsedResp.Error)
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
		return genderizeResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp genderizeResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return genderizeResponse{}, err
	}
	if parsedResp.Error != "" {
		return genderizeResponse{}, errors.New(parsedResp.Error)
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
		return nationalizeResponse{}, err
	}
	defer resp.Body.Close()

	var parsedResp nationalizeResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return nationalizeResponse{}, err
	}
	if parsedResp.Error != "" {
		return nationalizeResponse{}, errors.New(parsedResp.Error)
	}
	return parsedResp, nil
}
