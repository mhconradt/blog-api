package main


import (
	"fmt"
	"github.com/mattbaird/jsonpatch"
)

func main() {
	op := "replace"
	path := "/spec/containers/0/image"
	val := "gcr.io/$PROJECT_ID/blog:${TAG_NAME}"
	patch := jsonpatch.NewPatch(op, path, val)
	fmt.Println(patch)
}
