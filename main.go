package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== Bifrost Engine ===")
	fmt.Println("A modern game engine built in Go")
	fmt.Println()
	
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	
	command := os.Args[1]
	
	switch command {
	case "editor":
		fmt.Println("Starting Bifrost Engine Editor...")
		// TODO: Launch editor when ready
		fmt.Println("Editor not yet implemented. Use demos for now.")
		
	case "build":
		fmt.Println("Building project...")
		// TODO: Implement project building
		fmt.Println("Build system not yet implemented.")
		
	case "run":
		fmt.Println("Running project...")
		// TODO: Implement project running
		fmt.Println("Project runner not yet implemented.")
		
	case "demos":
		listDemos()
		
	case "--version", "-v":
		fmt.Println("Bifrost Engine v0.1.0")
		
	case "--help", "-h":
		printUsage()
		
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: bifrost <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  editor    Launch the Bifrost Engine editor")
	fmt.Println("  build     Build the current project")
	fmt.Println("  run       Run the current project")
	fmt.Println("  demos     List available demos")
	fmt.Println("  --version Show version information")
	fmt.Println("  --help    Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  bifrost editor")
	fmt.Println("  bifrost demos")
}

func listDemos() {
	fmt.Println("Available demos:")
	fmt.Println("  cd demos && go run basic_editor.go     - Interactive 3D editor")
	fmt.Println("  cd demos && go run shape_demo.go       - 3D shape viewer with controls")
	fmt.Println("  cd demos && go run cube_demo.go        - Rotating 3D cube")
	fmt.Println("  cd demos && go run camera_demo.go      - Camera orbiting demo")
	fmt.Println("  cd demos && go run main_demo.go        - Basic 3D cube display")
	fmt.Println()
	fmt.Println("Recommended: Start with 'basic_editor.go' for the full experience!")
}