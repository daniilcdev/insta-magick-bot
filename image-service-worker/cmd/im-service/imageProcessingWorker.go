package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/daniilcdev/insta-magick-bot/image-service-worker/config"
	types "github.com/daniilcdev/insta-magick-bot/image-service-worker/pkg"
)

type imageProcessingWorker struct {
	inDir  string
	outDir string
}

func NewProcessor(cfg *config.WorkerConfig) *imageProcessingWorker {
	return &imageProcessingWorker{
		inDir:  cfg.In,
		outDir: cfg.Out,
	}
}

func (im *imageProcessingWorker) Do(work types.Work) error {
	return im.doNow(&work)
}

func (im *imageProcessingWorker) doNow(work *types.Work) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()

	if work.Instruction == "" {
		return errors.New("no instruction")
	}

	inFile := im.inDir + work.File
	if err := saveImage(work.URL, inFile); err != nil {
		return err
	}

	if _, err = os.Stat(inFile); err != nil {
		return errors.New("image not found")
	}

	outFile := im.outDir + work.File

	args := strings.Split(inFile+" "+string(work.Instruction), " ")
	args = append(args, outFile)

	log.Printf("processing with filter '%s'\n", work.Filter)
	cmd := exec.Command("convert", args...)
	_, err = cmd.Output()

	return err
}

func saveImage(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non-200 response code")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
