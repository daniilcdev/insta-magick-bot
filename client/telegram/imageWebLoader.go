package telegram

import (
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
)

// TODO: implement DB storage
var Mu sync.Mutex = sync.Mutex{}
var ImgToChatMap map[string]string = make(map[string]string)

type imageWebLoader struct {
}

type downloadParams struct {
	url         string
	outFilename string
	outDir      string
	requesterId string
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

	defer Mu.Unlock()
	Mu.Lock()
	ImgToChatMap[params.outFilename] = params.requesterId

	return nil
}
