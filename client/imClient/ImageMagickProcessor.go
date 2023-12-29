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

func (im *IMProcessor) Naturalize(file fs.FileInfo) {
	// convert ./res/raw/*.jpg -normalize -auto-gamma ./res/processed/*.jpg
	fileName := file.Name()
	cmd := exec.Command("convert",
		im.inDir+fileName,
		"-normalize",
		"-auto-gamma",
		im.outDir+fileName,
	)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	// to debug
	fmt.Println(string(stdout))
}

func (im *IMProcessor) ProcessNewFile(path string, entry fs.DirEntry) {
	source := path + entry.Name()
	pending := im.inDir + entry.Name()

	err := os.Rename(source, pending)
	if err != nil {
		fmt.Printf("failed to move file %s: %v\n", entry.Name(), err)
		return
	}

	fi, _ := entry.Info()
	im.Naturalize(fi)

	fmt.Println("new file processed")
}
