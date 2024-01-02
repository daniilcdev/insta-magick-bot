package imclient

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

// const cmd_NORMALIZE_fmt = "%s -normalize -auto-gamma %s"

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

func (im *IMProcessor) Naturalize(filename string) {
	cmd := exec.Command("convert",
		im.inDir+filename,
		"-adaptive-sharpen",
		"5%",
		"-separate",
		"-contrast-stretch",
		"0.5%x0.5%",
		"-combine",
		"-enhance",
		"-auto-level",
		im.outDir+filename,
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
	filename := entry.Name()
	source := path + filename
	pending := im.inDir + filename

	err := os.Rename(source, pending)
	if err != nil {
		fmt.Printf("failed to move file %s: %v\n", filename, err)
		return
	}

	im.Naturalize(filename)
	os.Remove(pending)
}
