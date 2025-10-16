package main

import (
	"fmt"

	"github.com/dmikushin/podman-shared/version"
)

func main() {
	fmt.Print(version.Version.String())
}
