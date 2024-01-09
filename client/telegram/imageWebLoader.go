package telegram

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type imageWebLoader struct {
	storage Storage
}

type downloadParams struct {
	url         string
	outFilename string
	outDir      string
	requesterId string
	filter      string
}

func (iwl *imageWebLoader) downloadPhoto(params downloadParams) error {
	response, err := http.Get(params.url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non-200 response code")
	}

	file, err := os.Create(params.outDir + params.outFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	if params.filter == "" {
		params.filter = "Bright Summer"
	}

	filter, err := iwl.storage.FindFilter(params.filter)

	if err != nil {
		fmt.Println("filter not found")
		os.Remove(params.outDir + params.outFilename)
		return err
	}

	iwl.storage.CreateRequest(&NewRequest{
		File:        params.outFilename,
		RequesterId: params.requesterId,
		Filter:      filter.Name,
	})

	return nil
}
