package imclient

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

type IMProcessor struct {
	outDir string
}

func NewProcessor(outDir string) *IMProcessor {
	return &IMProcessor{
		outDir: outDir,
	}
}

// V1: mogrify -adaptive-sharpen 10% -separate -contrast-stretch 0.5%x0.5% -combine -enhance -auto-level -path %im.outDir %inDir/*.jpg
// V2: mogrify -adaptive-sharpen 10% -channel B -evaluate add 1.31 -channel G -evaluate add 1.37 +channel -modulate 120,142 -contrast-stretch -13%x-17% -enhance -path ../out *.jpg

func (im *IMProcessor) Beautify(inDir string) {
	cmd := exec.Command("mogrify",
		"-adaptive-sharpen", "10%",
		"-channel", "B", "-evaluate", "add", "1.31",
		"-channel", "G", "-evaluate", "add", "1.37",
		"+channel",
		"-modulate", "120,142",
		"-contrast-stretch", "-13%x-17%",
		"-enhance",
		"-path", im.outDir,
		inDir+"*.jpg",
	)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(stdout) > 0 {
		// to debug
		fmt.Println(string(stdout))
	}
}

func (im *IMProcessor) ProcessNewFilesInDir(path string, entries []fs.DirEntry) {
	im.Beautify(path)

	for _, entry := range entries {
		filename := entry.Name()
		pending := path + filename
		os.Remove(pending)
	}
}
