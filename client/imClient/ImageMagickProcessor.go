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

func (im *IMProcessor) Naturalize(inDir string) {
	cmd := exec.Command("mogrify",
		"-adaptive-sharpen",
		"10%",
		"-separate",
		"-contrast-stretch",
		"0.5%x0.5%",
		"-combine",
		"-enhance",
		"-auto-level",
		"-path",
		im.outDir,
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

func (im *IMProcessor) ProcessNewFile(path string, entries []fs.DirEntry) {
	im.Naturalize(path)

	for _, entry := range entries {
		filename := entry.Name()
		pending := path + filename
		os.Remove(pending)
	}
}
