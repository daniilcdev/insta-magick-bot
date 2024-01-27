package main

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/daniilcdev/insta-magick-bot/workers/im-worker/ports"
	"github.com/daniilcdev/insta-magick-bot/workers/im-worker/types"
)

type IMProcessor struct {
	outDir       string
	workingDir   string
	inDir        string
	workReporter ports.WorkReporter
}

func NewProcessor(cfg IMConfig) *IMProcessor {
	return &IMProcessor{
		inDir:      cfg.InDir(),
		outDir:     cfg.OutDir(),
		workingDir: "./res/tmp/",
	}
}

func (im *IMProcessor) WithWorkReporter(reporter ports.WorkReporter) *IMProcessor {
	im.workReporter = reporter
	return im
}

func (im *IMProcessor) OnWorkReceived(work types.Work) {
	err := im.applyFilter(&work)

	switch {
	case err != nil:
		log.Printf("IM process failed, error: '%v'\n", err)
		im.workReporter.Failed(work)
	default:
		im.workReporter.Done(work)
	}
}

func (im *IMProcessor) applyFilter(work *types.Work) error {
	if work.Filter == "" {
		return errors.New("no instruction")
	}

	log.Printf("[IMProcessor] processing files with filter %s\n", work.Filter)

	inFile := im.inDir + work.File
	outFile := im.outDir + work.File

	args := strings.Split(inFile+" "+string(work.Filter), " ")
	args = append(args, outFile)
	cmd := exec.Command("convert", args...)
	_, err := cmd.Output()

	switch errType := err.(type) {
	case *exec.ExitError:
		log.Printf("IM process failed: %s\n", string(errType.Stderr))
	case error:
		// do nothing there
	default:
		log.Printf("IM process exited unexpectedly: %s\n", errType.Error())
	}

	return err
}
