package hasher

import (
	"fmt"
	"testing"
)

func getTestImagePath(fileName string) string {
	return fmt.Sprintf("../test/images/%s", fileName)
}

func TestHashCalculation(t *testing.T) {
	hasher := New()
	hash, err := hasher.CalculateHash(getTestImagePath("same/1.jpg"))
	if err != nil {
		t.Error(err)
	}
	if hash.String() != "p:2adda34eaa5e41d1" {
		t.Error("hash not correct")
	}
}

func TestHashCalculationForUnexistingFile(t *testing.T) {
	hasher := New()
	hash, err := hasher.CalculateHash("unknown")
	if err == nil {
		t.Error("Expecting error for unexisting file")
	}
	if hash != nil {
		t.Error("Expecting no hash for unexisting file")
	}
}

// TODO: Make sure these images result in no difference
func _TestHashDifference(t *testing.T) {
	hasher := New()
	hash1, _ := hasher.CalculateHash(getTestImagePath("same/1.jpg"))
	hash2, _ := hasher.CalculateHash(getTestImagePath("same/2.jpg"))
	hash3, _ := hasher.CalculateHash(getTestImagePath("same/3.jpg"))
	hash4, _ := hasher.CalculateHash(getTestImagePath("same/4.jpg"))

	if distance, _ := hasher.CalculateDistance(hash1, hash2); distance != 0 {
		t.Errorf("Distance between hash 1 and 2 is %d", distance)
	}
	if distance, _ := hasher.CalculateDistance(hash2, hash3); distance != 0 {
		t.Errorf("Distance between hash 2 and 3 is %d", distance)
	}
	if distance, _ := hasher.CalculateDistance(hash3, hash4); distance != 0 {
		t.Errorf("Distance between hash 3 and 4 is %d", distance)
	}
}
