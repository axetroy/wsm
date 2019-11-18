package util

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/axetroy/terminal/internal/library/dotenv"
	"github.com/fatih/color"
)

func init() {
	err := dotenv.Load()

	if err != nil {
		log.Panicln(err)
	}
}

func PrintEnv() {
	envs := os.Environ()

	fmt.Println(color.GreenString("=== Runtime ==="))

	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Go OS: %s\n", runtime.GOOS)
	fmt.Printf("Go Arch: %s\n", runtime.GOARCH)

	fmt.Println(color.GreenString("=== Environmental Variable ==="))

	for _, e := range envs {
		fmt.Println(e)
	}
}
