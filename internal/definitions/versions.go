//go build -ldflags="-X 'main.version=2.2.0'"

package definitions

import (
	"flag"
	"fmt"
	"os"
)

func Version_handling(version string) {
	// flag name, default value, description
	versionFlag := flag.Bool("version", false, "Show version information")

	// Parse the flags
	flag.Parse()

	// If the version flag is set, print the version and exit
	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
}
