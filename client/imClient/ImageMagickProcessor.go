package imclient

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

type IMProcessor struct {
	inDir  string
	outDir string
}

func NewProcessor(sourceDir, outDir string) *IMProcessor {
	return &IMProcessor{
		inDir:  sourceDir,
		outDir: outDir,
	}
}

func (im *IMProcessor) Naturalize() {
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
		im.inDir+"*.jpg",
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

func (im *IMProcessor) ProcessNewFile(path string, entry fs.DirEntry) {
	im.Naturalize()

	filename := entry.Name()
	pending := im.inDir + filename
	os.Remove(pending)
}
