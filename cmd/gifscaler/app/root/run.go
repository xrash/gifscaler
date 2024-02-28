package root

import (
	"fmt"
	"os"
	"path/filepath"
)

func run(opts cliopts, args []string) int {
	fmt.Println(args)

	if len(args) < 1 {
		fmt.Println("There must be at least one arg")
		return 1
	}

	for _, arg := range args {
		err := do(arg, opts.keepWorkdir)

		if err == nil {
			fmt.Printf("resized %v\n", arg)
		} else {
			fmt.Printf("error resizing %v: %v\n", arg, err)
		}
	}

	return 0
}

func do(filename string, keepWorkdir bool) error {
	workdir, framesOriginalDir, framesScaledDir, err := createWorkdir()

	if !keepWorkdir {
		defer removeWorkdir(workdir)
	}

	if err != nil {
		return fmt.Errorf("error creating workdir: %w", err)
	}

	if err := saveFrames(filename, framesOriginalDir); err != nil {
		return fmt.Errorf("error saving frames: %w", err)
	}

	if err := scaleFrames(framesOriginalDir, framesScaledDir); err != nil {
		return fmt.Errorf("error scaling frames: %w", err)
	}

	paletteFilename := makePaletteFilename(workdir)

	if err := savePalette(framesScaledDir, paletteFilename); err != nil {
		return fmt.Errorf("error saving palette: %w", err)
	}

	outputFilename := makeOutputFilename(filename)

	if err := assembleOutput(framesScaledDir, paletteFilename, outputFilename); err != nil {
		return fmt.Errorf("error assembling output: %w", err)
	}

	return nil
}

func makePaletteFilename(workdir string) string {
	return filepath.Join(workdir, "palette.png")
}

func makeOutputFilename(filename string) string {
	return filename[:len(filename)-4] + "-scaled.gif"
}

func createWorkdir() (string, string, string, error) {
	workdir, err := os.MkdirTemp(os.TempDir(), "__gifscaler__")
	if err != nil {
		return "", "", "", fmt.Errorf("error from MkdirTemp: %w", err)
	}

	subdirs := []string{
		filepath.Join(workdir, "frames-original"),
		filepath.Join(workdir, "frames-scaled"),
	}

	for _, subdir := range subdirs {
		os.Mkdir(subdir, os.ModePerm)
	}

	return workdir, subdirs[0], subdirs[1], nil
}

func removeWorkdir(workdir string) {
	if err := os.RemoveAll(workdir); err != nil {
		fmt.Println("error removing workdir %v: %w", workdir, err)
	}
}
