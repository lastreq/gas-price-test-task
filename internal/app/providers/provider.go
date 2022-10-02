package providers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lastreq/gas-price-test-task/internal/app/model"
)

const (
	apiMock = "http://localhost:8081"
)

type Provider struct {
}

func New() Provider {
	prv := Provider{}
	return prv
}

func (prv Provider) GetGasHistory() (model.Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var data model.Data

	client := new(http.Client)

	request, err := http.NewRequestWithContext(ctx, "GET", apiMock, nil)
	if err != nil {
		return data, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return data, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
