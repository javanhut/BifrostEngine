package main

import (
	"fmt"
	"os"
	"os/exec"
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
		cmd := exec.Command("go", "run", "demos/gui_overlay_editor.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running editor: %v\n", err)
		}
		
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
		
	case "demo":
		if len(os.Args) < 3 {
			fmt.Println("Please specify a demo name")
			fmt.Println("Usage: bifrost demo <demo-name>")
			fmt.Println("Available demos: basic_editor, ui_editor, shape_demo, cube_demo, camera_demo")
			return
		}
		runDemo(os.Args[2])
		
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
	fmt.Println("Usage: bifrost <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  editor       Launch the Bifrost Engine UI editor")
	fmt.Println("  build        Build the current project")
	fmt.Println("  run          Run the current project")
	fmt.Println("  demos        List available demos")
	fmt.Println("  demo <name>  Run a specific demo")
	fmt.Println("  --version    Show version information")
	fmt.Println("  --help       Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  bifrost editor")
	fmt.Println("  bifrost demo ui_editor")
	fmt.Println("  bifrost demos")
}

func listDemos() {
	fmt.Println("Available demos:")
	fmt.Println("  overlay_editor   - Editor with on-screen GUI overlay")
	fmt.Println("  gui_editor       - GUI editor with menu-based controls")
	fmt.Println("  simple_ui_editor - Editor with keyboard/mouse controls")
	fmt.Println("  basic_editor     - Basic 3D editor with keyboard controls")
	fmt.Println("  shape_demo       - 3D shape viewer with controls")
	fmt.Println("  cube_demo        - Rotating 3D cube")
	fmt.Println("  camera_demo      - Camera orbiting demo")
	fmt.Println("  main_demo        - Basic 3D cube display")
	fmt.Println()
	fmt.Println("To run a demo: bifrost demo <demo-name>")
	fmt.Println("Example: bifrost demo overlay_editor")
	fmt.Println()
	fmt.Println("Recommended: Start with 'bifrost editor' for the on-screen GUI!")
}

func runDemo(demoName string) {
	var demoFile string
	switch demoName {
	case "overlay_editor":
		demoFile = "overlay_editor.go"
	case "gui_editor":
		demoFile = "gui_editor.go"
	case "simple_ui_editor":
		demoFile = "simple_ui_editor.go"
	case "ui_editor":
		demoFile = "ui_editor.go"
	case "basic_editor":
		demoFile = "basic_editor.go"
	case "shape_demo":
		demoFile = "shape_demo.go"
	case "cube_demo":
		demoFile = "cube_demo.go"
	case "camera_demo":
		demoFile = "camera_demo.go"
	case "main_demo":
		demoFile = "main_demo.go"
	default:
		fmt.Printf("Unknown demo: %s\n", demoName)
		fmt.Println("Run 'bifrost demos' to see available demos")
		return
	}
	
	fmt.Printf("Running %s demo...\n", demoName)
	cmd := exec.Command("go", "run", fmt.Sprintf("demos/%s", demoFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running demo: %v\n", err)
	}
}