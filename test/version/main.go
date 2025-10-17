package main

import (
	"fmt"

	"github.com/dmikushin/podman-shared/v5/version"
)

func main() {
	fmt.Print(version.Version.String())
}
