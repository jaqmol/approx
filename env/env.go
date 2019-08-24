package env

import (
	"log"
	"os"
	"strings"

	"github.com/jaqmol/approx/definition"
)

// AugmentMissing ...
func AugmentMissing(ds []definition.Definition) {
	missing := make([]string, 0)
	for _, d := range ds {
		for name, value := range d.Env {
			if value == nil {
				value := os.Getenv(name)
				if len(value) > 0 {
					d.Env[name] = &value
				} else {
					missing = append(missing, name)
				}
			}
		}
	}
	if len(missing) > 0 {
		log.Fatalf("Environment variables are missing: %v", strings.Join(missing, ", "))
	}
}
