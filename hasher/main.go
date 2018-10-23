package hasher

import (
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
	internalHash, err := goimagehash.PerceptionHash(img)
	hash = &ImageHash{internalHash: internalHash, filePath: filePath}
	return
}

// CalculateDistance  Calculate the distance between hashes
func (h *Hasher) CalculateDistance(hash1 *ImageHash, hash2 *ImageHash) (int, error) {
	return hash1.internalHash.Distance(hash2.internalHash)
}
