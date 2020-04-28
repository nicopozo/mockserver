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
	hash.Write([]byte(s)) //nolint:errcheck

	return hash.Sum32()
}
