package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"flag"
	"regexp"
	"path"
)

type Values struct {
	forSize     image.Point
	count       uint32
	// uint32 enough space for 16777216 images
	red_total   []uint32
	green_total []uint32
	blue_total  []uint32
}

func main() {
	folderPtr := flag.String("folder", "", "A folder with images")
	outputPtr := flag.String("output", "", "Output image (is always a jpg)")
	widthPtr := flag.Int("width", 0, "the expected with")
	heightPtr := flag.Int("height", 0, "the expected height")
	maxImagesPtr := flag.Int("maxImages", -1, "max amount of images to read, -1 for no max value")
	filenamePatternPtr := flag.String("filePattern", `(?i)\.jpe?g$`, "A reg exp pattern that must match the file.")

	flag.Parse()

	if *folderPtr == "" || *outputPtr == "" || *widthPtr == 0|| *heightPtr == 0 {
		flag.Usage()
		return
	}

	if *filenamePatternPtr == "" {
		*filenamePatternPtr = `.*`
	}

	pattern := regexp.MustCompile(*filenamePatternPtr)
	size := image.Point{*widthPtr, *heightPtr}
	execute(*folderPtr, pattern, *outputPtr, size, *maxImagesPtr)
}

func execute(folder string, filenamePattern *regexp.Regexp, outputImage string, size image.Point, imageCountToRead int) {
	files, err := ioutil.ReadDir(folder)
	handleError(err)

	xy := size.X * size.Y
	value := Values{size, 0, make([]uint32, xy), make([]uint32, xy), make([]uint32, xy)}

	i := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !filenamePattern.MatchString(file.Name()) {
			log.Printf("Ignore %s", file.Name())
			continue
		}

		log.Printf("Reading image %s", file.Name())
		if readImg(&value, path.Join(folder, file.Name())) {
			i++
		}

		if imageCountToRead > 0 && i >= imageCountToRead {
			break
		}
	}

	log.Printf("Read %d images", i)
	writeAverageImage(&value, outputImage)
}

func readImg(values *Values, imageFile string) bool {
	file, err := os.Open(imageFile)
	handleError(err)
	defer file.Close()
	image, _, err := image.Decode(file)

	size := image.Bounds().Size()
	if !size.Eq(values.forSize) {
		log.Printf("Image '%s' with %d x %d doesn't have the correct size (%d x %d)",
			imageFile, size.X, size.Y, values.forSize.X, values.forSize.Y)
		return false
	}

	values.count++

	var arrayCounter uint32 = 0
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, _ := image.At(x, y).RGBA()
			values.red_total[arrayCounter] += r
			values.green_total[arrayCounter] += g
			values.blue_total[arrayCounter] += b
			arrayCounter++
		}
	}

	return true
}

func writeAverageImage(values *Values, resultImage string) {
	if values.count == 0 {
		log.Println("No data for result image")
		return
	}

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, values.forSize})

	var arrayCounter uint32 = 0
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			r := uint8(values.red_total[arrayCounter] / values.count / 256)
			g := uint8(values.green_total[arrayCounter] / values.count / 256)
			b := uint8(values.blue_total[arrayCounter] / values.count / 256)
			img.Set(x, y, color.RGBA{r, g, b, 0xff})
			arrayCounter++
		}
	}

	file, err := os.Create(resultImage)
	handleError(err)
	defer file.Close()

	jpeg.Encode(file, img, &jpeg.Options{95})
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
