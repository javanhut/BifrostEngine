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
		args := []string{"run", "demos/gui_overlay_editor.go"}
		
		// Check for debug flag
		if len(os.Args) >= 3 && os.Args[2] == "--debug" {
			args = append(args, "--debug")
		}
		
		cmd := exec.Command("go", args...)
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
			fmt.Println("Usage: bifrost_engine demo <demo-name>")
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
	fmt.Println("Usage: bifrost_engine <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  editor [--debug]  Launch the Bifrost Engine UI editor")
	fmt.Println("  build             Build the current project")
	fmt.Println("  run               Run the current project")
	fmt.Println("  demos             List available demos")
	fmt.Println("  demo <name>       Run a specific demo")
	fmt.Println("  --version         Show version information")
	fmt.Println("  --help            Show this help message")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --debug           Enable detailed debug output for rendering")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  bifrost_engine editor")
	fmt.Println("  bifrost_engine editor --debug")
	fmt.Println("  bifrost_engine demo ui_editor")
	fmt.Println("  bifrost_engine demos")
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
	fmt.Println("🌟 LIGHTING DEMOS:")
	fmt.Println("  lighting_demo    - Basic lighting system demonstration")
	fmt.Println("  day_night_cycle  - Sun moving across sky with day/night")
	fmt.Println("  dramatic_lighting- Multiple colored lights with effects")
	fmt.Println("  studio_lighting  - Professional 3-point lighting setup")
	fmt.Println()
	fmt.Println("🎯 ASSET LOADING DEMOS:")
	fmt.Println("  obj_loading_demo - OBJ file loading and rendering demonstration")
	fmt.Println("  asset_browser_test - GUI asset browser functionality test")
	fmt.Println("  asset_buttons_test - Test asset browser button functionality")
	fmt.Println()
	fmt.Println("🎛️ GIZMO DEMOS:")
	fmt.Println("  gizmo_demo - Transform gizmo system demonstration")
	fmt.Println()
	fmt.Println("To run a demo: bifrost_engine demo <demo-name>")
	fmt.Println("Example: bifrost_engine demo day_night_cycle")
	fmt.Println()
	fmt.Println("Recommended: Start with 'bifrost_engine editor' for the on-screen GUI!")
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
	case "lighting_demo":
		demoFile = "lighting_demo.go"
	case "day_night_cycle":
		demoFile = "day_night_cycle.go"
	case "dramatic_lighting":
		demoFile = "dramatic_lighting.go"
	case "studio_lighting":
		demoFile = "studio_lighting.go"
	case "obj_loading_demo":
		demoFile = "obj_loading_demo.go"
	case "asset_browser_test":
		demoFile = "asset_browser_demo.go"
	case "asset_buttons_test":
		demoFile = "asset_buttons_demo.go"
	case "gizmo_demo":
		demoFile = "gizmo_demo.go"
	default:
		fmt.Printf("Unknown demo: %s\n", demoName)
		fmt.Println("Run 'bifrost_engine demos' to see available demos")
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