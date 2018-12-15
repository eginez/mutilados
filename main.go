package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

type mat [][]float64

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	imageFile, _ := os.Open("mutilados.png")
	defer imageFile.Close()

	img, err := png.Decode(imageFile)
	panicIfError(err)
	parts, err := generateParts(img)
	panicIfError(err)
	log.Println(len(parts))
}

func generateParts(fullImage image.Image) (parts []*image.RGBA, err error) {
	rgbImage := fullImage.(*image.RGBA)
	xParts := 3
	yParts := 2
	// the master image is not perfectly square hence the adjustment factors
	xPartSize := (fullImage.Bounds().Max.X / xParts) + 4
	yPartSize := (fullImage.Bounds().Max.Y / yParts) - 15
	parts = make([]*image.RGBA, 0)

	for y := 0; y < yParts; y++ {
		for x := 0; x < xParts; x++ {
			xCompensation := 0
			xFinalComp := 0
			compensation := 18
			// more tweaking to get account for not straight images
			if x == 1 && y == 1 {
				xCompensation = compensation
			}
			if x == 2 && y == 1 {
				xFinalComp = compensation
			}
			subImgRect := image.Rect(
				x*xPartSize+xFinalComp,
				y*yPartSize,
				(x+1)*xPartSize+xCompensation,
				(y+1)*yPartSize)
			log.Println(subImgRect)
			subImage := rgbImage.SubImage(subImgRect)
			parts = append(parts, subImage.(*image.RGBA))
			//err = saveImage(subImage, fmt.Sprintf("mutilados%d%d.png", x, y))
		}
	}
	return
}

func saveImage(img image.Image, name string) error {
	f, err := os.Create(name)
	defer f.Close()

	if err != nil {
		return err
	}

	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}

type myimg interface {
	At(x, y int) color.Color
	Bounds() image.Rectangle
}

func rotate(angle float64, img myimg) (res *image.RGBA, err error) {
	rotMatrix := rotationMarix(angle)
	posVector := newMat(1, 2)
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	res = image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.Y, img.Bounds().Max.X))
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			posVector.SetAt(0, 0, x)
			posVector.SetAt(0, 1, y)
			rotated, err := rotMatrix.Mul(posVector)
			if err != nil {
				return res, err
			}
			color := img.At(x, y)

			newX := rotated.At(0, 0)
			if newX < 0 {
				newX += maxX
			}

			newY := rotated.At(0, 1)
			if newY < 0 {
				newY += maxY
			}
			res.Set(newX, newY, color)
		}
	}
	return res, err
}

func rotationMarix(angle float64) mat {
	r := newMat(2, 2)
	r[0][0] = math.Cos(angle)
	r[1][0] = math.Sin(angle)
	r[0][1] = -math.Sin(angle)
	r[1][1] = math.Cos(angle)
	return r
}

func newMat(col, row int) mat {
	mat := make(mat, col)
	for i := range mat {
		mat[i] = make([]float64, row)
	}
	return mat
}

func (m mat) SetAt(x, y int, value int) {
	m.SetfAt(x, y, float64(value))
}

func (m mat) SetfAt(x, y int, value float64) {
	m[x][y] = value
}

func (m mat) At(x, y int) int {
	return int(math.Round(m[x][y]))
}

func (m mat) Cols() int {
	return len(m)
}

func (m mat) Rows() int {
	return len(m[0])
}

func (m mat) Mul(n mat) (res mat, err error) {
	res = newMat(n.Cols(), m.Rows())
	for k := 0; k < n.Cols(); k++ {
		for i := 0; i < m.Rows(); i++ {
			r := 0.0
			for j := 0; j < m.Cols(); j++ {
				r = r + m[j][i]*n[k][j]
			}
			res.SetfAt(k, i, r)
		}
	}
	return
}

func (m mat) String() string {
	b := strings.Builder{}
	for j := 0; j < m.Rows(); j++ {
		for i := 0; i < m.Cols(); i++ {
			b.WriteString(fmt.Sprintf("%d ", int(m[i][j])))
		}
		b.WriteString("\n")
	}
	return b.String()
}
