package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getImageFromFilePath(filePath string) (image.Image, string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	image, imageType, err := image.Decode(f)
	return image, imageType, err
}

func writeImage(filePath string, myImage image.Image, imageType string, sourceDir string, index int) error {
	// Save the files as PNG, because png is cool and autofills remaining pixels with transparent
	// Big respect to PNG
	filename := strings.Replace(filePath, sourceDir+string(os.PathSeparator), "", 1)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = "dest" + string(os.PathSeparator) + strconv.Itoa(index) + "-" + filename + ".png"
	fmt.Println(filename)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("Writing image " + filename)
	err = png.Encode(f, myImage)

	return err
}

func main() {
	var files []string
	sourceDir := "source"
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file == sourceDir {
			continue
		}
		fmt.Println("processing file : " + file)

		// read file, get its width / height and calculate the number of height in width
		myImage, imageType, _ := getImageFromFilePath(file)

		bounds := myImage.Bounds()
		width := bounds.Max.X
		height := bounds.Max.Y
		repeatedHeight := math.Ceil(float64(width) / float64(height))
		if width < height {
			fmt.Println("vas te faire foutre sérieux")
			continue
		}

		childTargetWidth := int(height)
		childStartWidth := 0
		targetWidth := height
		expectedRect := image.Rectangle{image.Point{0, 0}, image.Point{height, height}}

		for i := 0; i < int(repeatedHeight); i++ {
			targetBounds := image.Rectangle{image.Point{childStartWidth, 0}, image.Point{int(childTargetWidth), height}}
			subImage := myImage.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(targetBounds)
			if subImage.Bounds() != targetBounds {
				mySubImage := image.NewRGBA(expectedRect)
				fakeX := 0
				for x := subImage.Bounds().Min.X; x < subImage.Bounds().Max.X; x++ {
					fakeY := 0
					for y := subImage.Bounds().Min.Y; y < subImage.Bounds().Max.Y; y++ {
						mySubImage.Set(fakeX, fakeY, subImage.At(x, y))
						fakeY++
					}
					fakeX++
				}
				subImage = mySubImage
			}
			childStartWidth = childTargetWidth + 1
			childTargetWidth += targetWidth + 1

			writeImage(file, subImage, imageType, sourceDir, i)
		}
	}
}
