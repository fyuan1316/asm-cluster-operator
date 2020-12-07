package lib

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"hash"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
)

func ComputeHash(obj interface{} /*, collisionCount *int32*/) string {
	templateHasher := fnv.New32a()
	DeepHashObject(templateHasher, obj)

	// Add collisionCount in the hash if it exists.
	//if collisionCount != nil {
	//	collisionCountBytes := make([]byte, 8)
	//	binary.LittleEndian.PutUint32(collisionCountBytes, uint32(*collisionCount))
	//	templateHasher.Write(collisionCountBytes)
	//}
	return rand.SafeEncodeString(fmt.Sprint(templateHasher.Sum32()))
}

func DeepHashObject(hasher hash.Hash, objectToWrite interface{}) {
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", objectToWrite)
}
