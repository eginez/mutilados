package main

import (
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultOne(t *testing.T) {
	a := newMat(1, 2)
	a.SetAt(0, 0, 1)
	a.SetAt(0, 1, 1)

	b := newMat(2, 2)
	b.SetAt(0, 0, 1)
	b.SetAt(0, 1, 1)
	b.SetAt(1, 0, 1)
	b.SetAt(1, 1, 1)

	c, err := b.Mul(a)
	assert.NoError(t, err)
	assert.Equal(t, 2, c.At(0, 0))
	assert.Equal(t, 2, c.At(0, 1))
}

func TestRotation(t *testing.T) {
	imageFile, _ := os.Open("out/mutilados00.png")
	defer imageFile.Close()

	img, _ := png.Decode(imageFile)

	angle := math.Pi / 2
	r1, _ := rotate(angle, img)
	r2, _ := rotate(2*angle, img)
	r3, _ := rotate(3*angle, img)
	r4, _ := rotate(4*angle, img)

	saveImage(r1, "out/r1.png")
	saveImage(r2, "out/r2.png")
	saveImage(r3, "out/r3.png")
	saveImage(r4, "out/r4.png")

}
