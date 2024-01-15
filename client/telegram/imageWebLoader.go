package telegram

import (
	"errors"
	"io"
	"net/http"
	"os"
)

type imageWebLoader struct {
	storage Storage
	outDir  string
}

type downloadParams struct {
	url         string
	outFilename string
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

	file, err := os.Create(iwl.outDir + params.outFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	iwl.storage.CreateRequest(&NewRequest{
		File:        params.outFilename,
		RequesterId: params.requesterId,
		Filter:      params.filter,
	})

	return nil
}
