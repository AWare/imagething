package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Println(len(os.Args))
		fmt.Println("Generates a png. \n Need to call with: path of image\n path of new image\n probability of block")
	}
	old := os.Args[1]
	new := os.Args[2]
	p, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	f, _ := os.Create(new)
	w := bufio.NewWriter(f)
	defer f.Close()
	defer w.Flush()

	sourceImage, err := getImageFromPath(old)
	if err != nil {
		log.Fatal(err)
	}

	resultImage, err := transformImage(sourceImage, p)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(w, resultImage)
}

func transformImage(img image.Image, p int) (image.Image, error) {
	output := image.NewRGBA(img.Bounds())
	rectangles := splitSquare(img.Bounds(), 60)
	draw.Draw(output, output.Bounds(), img, image.ZP, draw.Src)
	for _, r := range rectangles {
		if rand.Intn(p+1) == 0 {
			c, err := getAverageColorForRegion(img, r)
			if err != nil {
				return nil, err
			}
			draw.Draw(output, r, &image.Uniform{c}, image.ZP, draw.Src)
		}
	}
	return output, nil

}

func getImageFromPath(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	return img, err
}

func getImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(resp.Body)
	return img, err
}

func getAverageColorForRegion(img image.Image, rc image.Rectangle) (color.Color, error) {
	if !rc.In(img.Bounds()) {
		return nil, fmt.Errorf("r is not withing img bounds")
	}

	var red, green, blue, alpha uint
	pixels := uint(rc.Dx() * rc.Dy())
	for i := rc.Min.X; i < rc.Max.X; i++ {
		for j := rc.Min.Y; j < rc.Max.Y; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			red += uint(r)
			green += uint(g)
			blue += uint(b)
			alpha += uint(a)
		}
	}
	return color.RGBA64{uint16(red / pixels), uint16(green / pixels), uint16(blue / pixels), uint16(alpha / pixels)}, nil

}

func colorDistance(co1, co2 color.Color) int {
	//This returns the square of the color distance.
	m := color.YCbCrModel
	c2 := m.Convert(co2).(color.YCbCr)
	c1 := m.Convert(co1).(color.YCbCr)
	return int((c2.Y-c1.Y)*(c2.Y-c1.Y) + (c2.Cb-c1.Cb)*(c2.Cb-c1.Cb) + (c2.Cr-c1.Cr)*(c2.Cr-c1.Cr))
}

func splitSquare(r image.Rectangle, size int) []image.Rectangle {

	var nx, ny int
	nx = ((r.Max.X - r.Min.X) / size)
	ny = ((r.Max.Y - r.Min.Y) / size)
	//do all that 2d array crap
	rectangles := make([]image.Rectangle, 0, nx*ny)
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			x := r.Min.X + i*size
			y := r.Min.Y + j*size

			rect := image.Rect(x, y, x+size, y+size)

			rectangles = append(rectangles, rect)

		}
	}
	return rectangles
}
