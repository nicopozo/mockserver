package stringutils

import (
	"hash/fnv"
)

func ArraysContains(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}

func Hash(s string) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(s))

	return hash.Sum32()
}
