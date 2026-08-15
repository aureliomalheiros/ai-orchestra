package main

import (
	"log"
	"github.com/aureliomalheiros/ai-orchestra/internal/bubbletea"
)

func main() {
	if _, err := bubbletea.NewProgram(nil).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
