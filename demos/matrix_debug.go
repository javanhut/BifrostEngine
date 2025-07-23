package main

import (
	"fmt"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/scene"
)

func main() {
	// Create transform
	transform := scene.NewTransform()
	transform.SetPosition(bmath.NewVector3(0, 0, 0))
	
	// Get matrices
	localMatrix := transform.GetLocalMatrix()
	worldMatrix := transform.GetWorldMatrix()
	identity := bmath.NewMatrix4Identity()
	
	fmt.Printf("Local matrix: %v\n", localMatrix)
	fmt.Printf("World matrix: %v\n", worldMatrix)
	fmt.Printf("Identity matrix: %v\n", identity)
	fmt.Printf("Local == Identity: %v\n", localMatrix == identity)
	fmt.Printf("World == Identity: %v\n", worldMatrix == identity)
}