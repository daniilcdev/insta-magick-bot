package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/daniilcdev/insta-magick-bot/workers/im-worker/config"
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
)

type IMProcessor struct {
	outDir     string
	workingDir string
	inDir      string
}

func NewProcessor(cfg config.IMConfig) *IMProcessor {
	return &IMProcessor{
		inDir:      cfg.InDir(),
		outDir:     cfg.OutDir(),
		workingDir: "./res/tmp/",
	}
}

func (im *IMProcessor) Do(work types.Work) error {
	err := im.doNow(&work)
	return err
}

func (im *IMProcessor) doNow(work *types.Work) error {
	if work.Filter == "" {
		return errors.New("no instruction")
	}

	if _, err := os.Stat(im.inDir + work.File); err != nil {
		return errors.New("image not found")
	}

	log.Printf("[IMProcessor] processing files with filter %s\n", work.Filter)

	inFile := im.inDir + work.File
	outFile := im.outDir + work.File

	args := strings.Split(inFile+" "+string(work.Filter), " ")
	args = append(args, outFile)
	cmd := exec.Command("convert", args...)
	_, err := cmd.Output()

	return err
}
