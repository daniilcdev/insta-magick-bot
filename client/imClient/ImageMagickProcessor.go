package imclient

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/daniilcdev/insta-magick-bot/client/imClient/ports"
	"github.com/daniilcdev/insta-magick-bot/client/telegram"
)

type IMProcessor struct {
	outDir            string
	workingDir        string
	db                telegram.Storage
	completionHandler ports.CompletionHandler
}

func NewProcessor(cfg IMConfig, db telegram.Storage) *IMProcessor {
	return &IMProcessor{
		outDir:     cfg.OutDir(),
		workingDir: "./res/tmp/",
		db:         db,
	}
}

func (imc *IMProcessor) WithCompletionHandler(handler ports.CompletionHandler) *IMProcessor {
	imc.completionHandler = handler
	return imc
}

func (im *IMProcessor) Beautify(inDir string, filterName string, files []string) error {
	filter, err := im.db.FindFilter(filterName)

	if err != nil {
		log.Printf("filter %s not founds, %v", filterName, err)
		return err
	}

	log.Printf("[IMProcessor] processing files with filter %s\n", filter.Name)

	for _, file := range files {
		err := os.Rename(inDir+file, im.workingDir+file)
		if err != nil {
			log.Printf("unable to move file %s, %v", file, err)
		}
	}

	args := strings.Split(filter.Receipt, " ")
	args = append(args, "-path", im.outDir, im.workingDir+"*.jpg")
	cmd := exec.Command("mogrify", args...)
	_, err = cmd.Output()

	switch errType := err.(type) {
	case *exec.ExitError:
		log.Printf("IM process failed: %s\n", string(errType.Stderr))
	case nil:
		var misses int = 0
		for _, file := range files {
			err := os.Remove(im.workingDir + file)
			if err != nil {
				misses++
			}
		}

		if misses > 0 {
			log.Printf("[WARN] not all files were removed from temp dir\n")
		}

	default:
		log.Printf("IM process exited unexpectedly: %s\n", errType.Error())
	}

	return err
}

func (im *IMProcessor) ProcessNewFilesInDir(path string) {
	const batchSize = 10
	pending := im.db.Schedule(batchSize)

	if len(pending) == 0 {
		return
	}

	m := make(map[string][]string)
	for _, row := range pending {
		s, exists := m[row.FilterName]
		if exists {
			m[row.FilterName] = append(s, row.File)
			continue
		}

		s = []string{row.File}
		m[row.FilterName] = s
	}

	for filter, files := range m {
		err := im.Beautify(path, filter, files)

		switch err {
		case nil:
			im.db.CompleteRequests(files)
		default:
			log.Printf("IM process failed, rolling back; error: '%v'\n", err)
			im.db.Rollback(files)
		}
	}

	im.completionHandler.OnProcessCompleted(im.outDir)
}
