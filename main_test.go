package godeps

import (
	"fmt"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	var v, err = parseVersion("v1.2.3")
	fmt.Printf("%+v\n", v)
	fmt.Println(err)

	v, err = parseVersion("v0.0.0-20230522175609-2e198f4a06a1")
	fmt.Printf("%+v\n", v)
	fmt.Println(err)
}
