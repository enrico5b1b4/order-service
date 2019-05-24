package env

import (
	"fmt"
	"log"
	"os"
)

func MustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}
	return value
}
