package main

import (
	"image"
	_ "image/png" //needed to decode png
	"log"
	"math"
	"os"

	"github.com/fogleman/gg"
	"fmt"
)

const (
	WIDTH  = 1080
	HEIGHT = 1080
)

//https://jonnoftw.github.io/2017/01/18/markov-chain-image-generation
func main() {
	image, err := openImage("input/cat.png")

	if err != nil {
		log.Fatalf("Could not open image: %v", err)
	}

	fmt.Println("Making chain")
	markovChain := makeSecondOrderChain(image)

	context := gg.NewContext(WIDTH, HEIGHT)
	fmt.Println("Drawing image")
	drawImage(context, markovChain)

	context.SavePNG("output/markov.png")
}

func openImage(filename string) (image.Image, error) {
	imageReader, openErr := os.Open(filename)

	if openErr != nil {
		return nil, openErr
	}

	img, _, decodeErr := image.Decode(imageReader)

	if decodeErr != nil {
		return nil, decodeErr
	}

	return img, nil
}

func drawImage(context *gg.Context, chain *SecondOrderMarkovChain) {
	var x, y int

	for y = 0; y < WIDTH; y++ {
		for x = 0; x < HEIGHT; x++ {
			c := chain.Next()

			//fmt.Println(c)

			context.SetColor(c)
			context.SetPixel(x, y)
		}
	}
}

func makeChain(image image.Image) *MarkovChain {
	markovChain := NewChain()

	minY := image.Bounds().Min.Y
	minX := image.Bounds().Min.X
	maxY := image.Bounds().Max.Y
	maxX := image.Bounds().Max.X

	for col := minY; col < maxY; col++ {
		for row := minX; row < maxX; row++ {

			//add neighbors
			for dH := -1; dH <= 1; dH++ {
				for dV := -1; dV <= 1; dV++ {
					if dV != 0 && dH != 0 {
						markovChain.Add(
							image.At(clamp(row, minX, maxX-1), clamp(col, minY, maxY-1)),
							image.At(clamp(row+dV, minX, maxX-1), clamp(col+dH, minY, maxY-1)))
					}
				}
			}

		}
	}

	return markovChain
}

func makeSecondOrderChain(image image.Image) *SecondOrderMarkovChain {
	markovChain := NewSecondOrderChain()

	minY := image.Bounds().Min.Y
	minX := image.Bounds().Min.X
	maxY := image.Bounds().Max.Y
	maxX := image.Bounds().Max.X

	for col := minY; col < maxY; col++ {
		for row := minX; row < maxX; row++ {

			//add neighbors
			for dH := -1; dH <= 1; dH++ {
				for dV := -1; dV <= 1; dV++ {
					for dH_2 := -1; dH_2 <= 1; dH_2++ {
						for dV_2 := -1; dV_2 <= 1; dV_2++ {
							if dV != 0 && dH != 0 && dV_2 != 0 && dH_2 != 0{
								markovChain.Add(
									image.At(clamp(row, minX, maxX-1), clamp(col, minY, maxY-1)),
									image.At(clamp(row+dV, minX, maxX-1), clamp(col+dH, minY, maxY-1)),
									image.At(clamp(row+dV+dV_2, minX, maxX-1), clamp(col+dH+dH_2, minY, maxY-1)))
							}
						}
					}
				}
			}

		}
	}

	return markovChain
}

func clamp(x, min, max int) int {
	return int(math.Max(float64(min), math.Min(float64(max), float64(x))))
}
