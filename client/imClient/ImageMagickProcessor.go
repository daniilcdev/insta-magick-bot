package imclient

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/daniilcdev/insta-magick-bot/client/imClient/ports"
	"github.com/daniilcdev/insta-magick-bot/client/telegram"
)

type IMProcessor struct {
	outDir            string
	db                telegram.Storage
	completionHandler ports.CompletionHandler
}

func NewProcessor(outDir string, db telegram.Storage) *IMProcessor {
	return &IMProcessor{
		outDir: outDir,
		db:     db,
	}
}

func (imc *IMProcessor) WithCompletionHandler(handler ports.CompletionHandler) *IMProcessor {
	imc.completionHandler = handler
	return imc
}

// V1: mogrify -adaptive-sharpen 10% -separate -contrast-stretch 0.5%x0.5% -combine -enhance -auto-level -path %im.outDir %inDir/*.jpg
// V2: mogrify -adaptive-sharpen 10% -channel B -evaluate add 1.31 -channel G -evaluate add 1.37 +channel -modulate 120,142 -contrast-stretch -13%x-17% -enhance -path ../out *.jpg

func (im *IMProcessor) Beautify(inDir string, filterName string, files []string) {
	filter, err := im.db.FindFilter(filterName)

	if err != nil {
		log.Printf("filter %s not founds, %v", filterName, err)
		return
	}

	log.Printf("[IMProcessor] processing files with filter %s\n", filter.Name)

	const cwd = "./res/tmp/"

	for _, file := range files {
		err := os.Rename(inDir+file, cwd+file)
		if err != nil {
			log.Printf("unable to move file %s, %v", file, err)
		}
	}

	args := strings.Split(filter.Receipt, " ")
	args = append(args, "-path", im.outDir, cwd+"*.jpg")
	cmd := exec.Command("mogrify", args...)
	stdout, err := cmd.Output()

	for _, file := range files {
		err := os.Remove(cwd + file)
		if err != nil {
			log.Printf("unable to delete file %s, %v", file, err)
		}
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(stdout) > 0 {
		// to debug
		fmt.Println(string(stdout))
	}
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
		im.Beautify(path, filter, files)
		im.db.CompleteRequests(files)
	}

	im.completionHandler.OnProcessCompleted(im.outDir)
}
