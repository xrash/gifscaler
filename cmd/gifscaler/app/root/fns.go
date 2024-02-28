package root

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
)

func saveFrames(filename, destination string) error {
	// ffmpeg actually understands %d and will replace it with the frame number,
	// producing a bunch of files like frame_0, frame_1, frame_2 etc.
	cmd := exec.Command("ffmpeg", "-i", filename, destination+"/frame_%d.png")

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func scaleFrames(source, destination string) error {

	visit := func(path string, info fs.FileInfo, receivedErr error) error {
		// Skip the first path, which is the directory itself, instead of a file.
		if info.Name()[len(info.Name())-4:] != ".png" {
			return nil
		}

		if receivedErr != nil {
			return fmt.Errorf("visit fn received an error: %w", receivedErr)
		}

		scaledFilename := filepath.Join(
			destination,
			info.Name()[:len(info.Name())-4]+"-scaled"+".png",
		)

		//fmt.Println("path", path, info.Name(), scaledFilename)

		return scale(path, scaledFilename)
	}

	filepath.Walk(source, visit)

	return nil
}

func savePalette(framesDirname, paletteFilename string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		framesDirname+"/frame_%d-scaled.png",
		"-vf",
		"palettegen",
		paletteFilename,
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func assembleOutput(framesDirname, paletteFilename, outputFilename string) error {
	fmt.Println(framesDirname + "/frame_%d-scaled.png")
	fmt.Println(outputFilename)

	cmd := exec.Command(
		"ffmpeg",
		"-i",
		framesDirname+"/frame_%d-scaled.png",
		"-i",
		paletteFilename,
		"-lavfi",
		"paletteuse,setpts=6*PTS",
		"-r",
		"16",
		"-y",
		outputFilename,
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func scale(input, output string) error {
	// ffmpeg actually understands %d and will replace it with the frame number,
	// producing a bunch of files like frame_0, frame_1, frame_2 etc.
	cmd := exec.Command("scalex", "-k", "2", input, output)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
