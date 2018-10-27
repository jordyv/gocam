package hasher

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/corona10/goimagehash"
)

// Hasher  Utility to calculate hashes for images
type Hasher struct{}

// ImageHash  Represents the unique hash for an image
type ImageHash struct {
	internalHash *goimagehash.ImageHash
	filePath     string
}

func (i *ImageHash) String() string {
	return i.internalHash.ToString()
}

// New  Create new instance of an Hasher
func New() *Hasher {
	return &Hasher{}
}

func (h *Hasher) prepareImage(target image.Image) image.Image {
	size := target.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)
	// loop though all the x
	for x := 0; x < size.X; x++ {
		// and now loop thorough all of this x's y
		for y := 0; y < size.Y; y++ {
			pixel := target.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			r := float64(originalColor.R)
			g := float64(originalColor.G)
			b := float64(originalColor.B)
			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey, G: grey, B: grey, A: originalColor.A,
			}
			wImg.Set(x, y, c)
		}
	}
	return wImg
}

// CalculateHash  Calculate the unique hash for an image
func (h *Hasher) CalculateHash(filePath string) (hash *ImageHash, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		return
	}
	img = h.prepareImage(img)

	internalHash, err := goimagehash.PerceptionHash(img)
	hash = &ImageHash{internalHash: internalHash, filePath: filePath}
	return
}

// CalculateDistance  Calculate the distance between hashes
func (h *Hasher) CalculateDistance(hash1 *ImageHash, hash2 *ImageHash) (int, error) {
	return hash1.internalHash.Distance(hash2.internalHash)
}
