package stringutils

import (
	"hash/fnv"
)

func Hash(s string) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(s))

	return hash.Sum32()
}
