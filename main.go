package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func main() {
	iconPath := os.Getenv("stamp_path_to_icons")
	envName := os.Getenv("stamp_env_name")
	buildNumber := os.Getenv("stamp_build_number")
	topFgColor := os.Getenv("top_foreground_color")
	topBgColor := os.Getenv("top_background_color")
	bottomFgColor := os.Getenv("bottom_foreground_color")
	bottomBgColor := os.Getenv("bottom_background_color")

	fmt.Println("Build number to stamp:", buildNumber)
	fmt.Println("Environemtn name to stamp:", envName)
	fmt.Println("Top foreground color is:", topFgColor)
	fmt.Println("Top background color is:", topBgColor)
	fmt.Println("Bottom foreground color is:", bottomFgColor)
	fmt.Println("Bottom background color is:", bottomBgColor)

	fmt.Println("Finding icons from directory:", iconPath)

	files, err := ioutil.ReadDir(iconPath)
	if err != nil {
		fmt.Println("Could not read directory!")
		os.Exit(1)
	}

	filePaths := make([]string, 0)
	for _, f := range files {
		if path.Ext(f.Name()) == ".png" {
			filePaths = append(filePaths, path.Join(iconPath, f.Name()))
		}
	}

	fmt.Println(filePaths)

	for _, f := range filePaths {
		dimsOut, err := exec.Command("identify", "-format", "%w,%h", f).Output()
		if err != nil {
			fmt.Println("ImageMagick failed!")
			os.Exit(1)
		}

		dims := strings.Split(string(dimsOut), ",")

		width, err1 := strconv.Atoi(dims[0])
		height, err2 := strconv.Atoi(dims[1])
		if err1 != nil && err2 != nil {
			os.Exit(1)
		}

		bannerH := int(math.Floor(float64(height) * 0.25))
		bannerDims := strconv.Itoa(width) + "x" + strconv.Itoa(bannerH)

		imgOutString, error := exec.Command("convert", f,
			"-background", topBgColor,
			"-fill", topFgColor,
			"-gravity", "north",
			"-weight", "700",
			"-size", bannerDims,
			"caption:"+envName,
			"-composite",
			"(",
			"-background", bottomBgColor,
			"-fill", bottomFgColor,
			"-gravity", "south",
			"-weight", "700",
			"-size", bannerDims,
			"caption:"+buildNumber,
			")",
			"-composite", f).CombinedOutput()
		if error != nil {
			fmt.Println(string(imgOutString))
			fmt.Println(error)
			fmt.Println("ImageMagick failed!")
			os.Exit(1)
		}

	}

	//
	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}
