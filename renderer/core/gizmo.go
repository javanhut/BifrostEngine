package core

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/opengl"
)

type GizmoType int

const (
	GizmoTranslate GizmoType = iota
	GizmoRotate
	GizmoScale
)

type GizmoAxis int

const (
	GizmoAxisNone GizmoAxis = iota
	GizmoAxisX
	GizmoAxisY
	GizmoAxisZ
	GizmoAxisAll
)

type Gizmo struct {
	Type           GizmoType
	selectedAxis   GizmoAxis
	hoveredAxis    GizmoAxis
	lineShader     *opengl.Shader
	position       bmath.Vector3
	scale          float32
	visible        bool
	
	// Line meshes for each axis
	xAxisMesh      *opengl.Mesh
	yAxisMesh      *opengl.Mesh
	zAxisMesh      *opengl.Mesh
	
	// Additional meshes for specific gizmo types
	xyPlaneMesh    *opengl.Mesh // For translate
	xzPlaneMesh    *opengl.Mesh
	yzPlaneMesh    *opengl.Mesh
	
	// Colors
	xColor         bmath.Vector3
	yColor         bmath.Vector3
	zColor         bmath.Vector3
	selectedColor  bmath.Vector3
	hoveredColor   bmath.Vector3
}

func NewGizmo(lineShader *opengl.Shader) *Gizmo {
	g := &Gizmo{
		Type:          GizmoTranslate,
		selectedAxis:  GizmoAxisNone,
		hoveredAxis:   GizmoAxisNone,
		lineShader:    lineShader,
		scale:         1.0,
		visible:       true,
		xColor:        bmath.NewVector3(1.0, 0.0, 0.0), // Red
		yColor:        bmath.NewVector3(0.0, 1.0, 0.0), // Green
		zColor:        bmath.NewVector3(0.0, 0.0, 1.0), // Blue
		selectedColor: bmath.NewVector3(1.0, 1.0, 0.0), // Yellow
		hoveredColor:  bmath.NewVector3(1.0, 0.5, 0.0), // Orange
	}
	
	g.createMeshes()
	return g
}

func (g *Gizmo) createMeshes() {
	// Create translation gizmo meshes (arrows)
	g.createTranslationMeshes()
	
	// TODO: Create rotation gizmo meshes (circles)
	// TODO: Create scale gizmo meshes (boxes)
}

func (g *Gizmo) createTranslationMeshes() {
	// X-axis arrow (pointing right)
	xVertices := []float32{
		// Main line
		0.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		// Arrow head
		1.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.9, 0.1, 0.0, 1.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.9, -0.1, 0.0, 1.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.9, 0.0, 0.1, 1.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.9, 0.0, -0.1, 1.0, 0.0, 0.0,
	}
	xIndices := []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	g.xAxisMesh = opengl.NewLinesMesh(xVertices, xIndices)
	
	// Y-axis arrow (pointing up)
	yVertices := []float32{
		// Main line
		0.0, 0.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 0.0, 1.0, 0.0,
		// Arrow head
		0.0, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.1, 0.9, 0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 0.0, 1.0, 0.0,
		-0.1, 0.9, 0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.9, 0.1, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.9, -0.1, 0.0, 1.0, 0.0,
	}
	yIndices := []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	g.yAxisMesh = opengl.NewLinesMesh(yVertices, yIndices)
	
	// Z-axis arrow (pointing forward)
	zVertices := []float32{
		// Main line
		0.0, 0.0, 0.0, 0.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		// Arrow head
		0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.1, 0.0, 0.9, 0.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		-0.1, 0.0, 0.9, 0.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.0, 0.1, 0.9, 0.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.0, -0.1, 0.9, 0.0, 0.0, 1.0,
	}
	zIndices := []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	g.zAxisMesh = opengl.NewLinesMesh(zVertices, zIndices)
	
	// Create plane indicators for 2D translation
	g.createPlaneMeshes()
}

func (g *Gizmo) createPlaneMeshes() {
	// XY plane (small square)
	xyVertices := []float32{
		0.2, 0.2, 0.0, 1.0, 1.0, 0.0,
		0.4, 0.2, 0.0, 1.0, 1.0, 0.0,
		0.4, 0.2, 0.0, 1.0, 1.0, 0.0,
		0.4, 0.4, 0.0, 1.0, 1.0, 0.0,
		0.4, 0.4, 0.0, 1.0, 1.0, 0.0,
		0.2, 0.4, 0.0, 1.0, 1.0, 0.0,
		0.2, 0.4, 0.0, 1.0, 1.0, 0.0,
		0.2, 0.2, 0.0, 1.0, 1.0, 0.0,
	}
	xyIndices := []uint32{0, 1, 2, 3, 4, 5, 6, 7}
	g.xyPlaneMesh = opengl.NewLinesMesh(xyVertices, xyIndices)
	
	// XZ plane
	xzVertices := []float32{
		0.2, 0.0, 0.2, 1.0, 0.0, 1.0,
		0.4, 0.0, 0.2, 1.0, 0.0, 1.0,
		0.4, 0.0, 0.2, 1.0, 0.0, 1.0,
		0.4, 0.0, 0.4, 1.0, 0.0, 1.0,
		0.4, 0.0, 0.4, 1.0, 0.0, 1.0,
		0.2, 0.0, 0.4, 1.0, 0.0, 1.0,
		0.2, 0.0, 0.4, 1.0, 0.0, 1.0,
		0.2, 0.0, 0.2, 1.0, 0.0, 1.0,
	}
	xzIndices := []uint32{0, 1, 2, 3, 4, 5, 6, 7}
	g.xzPlaneMesh = opengl.NewLinesMesh(xzVertices, xzIndices)
	
	// YZ plane
	yzVertices := []float32{
		0.0, 0.2, 0.2, 0.0, 1.0, 1.0,
		0.0, 0.4, 0.2, 0.0, 1.0, 1.0,
		0.0, 0.4, 0.2, 0.0, 1.0, 1.0,
		0.0, 0.4, 0.4, 0.0, 1.0, 1.0,
		0.0, 0.4, 0.4, 0.0, 1.0, 1.0,
		0.0, 0.2, 0.4, 0.0, 1.0, 1.0,
		0.0, 0.2, 0.4, 0.0, 1.0, 1.0,
		0.0, 0.2, 0.2, 0.0, 1.0, 1.0,
	}
	yzIndices := []uint32{0, 1, 2, 3, 4, 5, 6, 7}
	g.yzPlaneMesh = opengl.NewLinesMesh(yzVertices, yzIndices)
}

func (g *Gizmo) SetPosition(pos bmath.Vector3) {
	g.position = pos
}

func (g *Gizmo) SetScale(scale float32) {
	g.scale = scale
}

func (g *Gizmo) SetVisible(visible bool) {
	g.visible = visible
}

func (g *Gizmo) SetType(gizmoType GizmoType) {
	g.Type = gizmoType
}

func (g *Gizmo) SetHoveredAxis(axis GizmoAxis) {
	g.hoveredAxis = axis
}

func (g *Gizmo) SetSelectedAxis(axis GizmoAxis) {
	g.selectedAxis = axis
}

func (g *Gizmo) GetSelectedAxis() GizmoAxis {
	return g.selectedAxis
}

// GetHoveredAxis returns the currently hovered axis
func (g *Gizmo) GetHoveredAxis() GizmoAxis {
	return g.hoveredAxis
}

// HandleMouseMove processes mouse movement for hover detection
func (g *Gizmo) HandleMouseMove(mouseX, mouseY float64, camera *bmath.Vector3, view, projection bmath.Matrix4, screenWidth, screenHeight int) {
	if !g.visible {
		g.hoveredAxis = GizmoAxisNone
		return
	}
	
	// Convert mouse coordinates to normalized device coordinates (-1 to 1)
	ndcX := (2.0 * mouseX / float64(screenWidth)) - 1.0
	ndcY := 1.0 - (2.0 * mouseY / float64(screenHeight))
	
	// Create ray from camera through mouse position
	ray := g.createMouseRay(ndcX, ndcY, camera, view, projection)
	
	// Test ray intersection with gizmo axes
	g.hoveredAxis = g.testRayIntersection(ray)
}

// HandleMouseClick processes mouse clicks for selection
func (g *Gizmo) HandleMouseClick(mouseX, mouseY float64, camera *bmath.Vector3, view, projection bmath.Matrix4, screenWidth, screenHeight int) GizmoAxis {
	if !g.visible {
		return GizmoAxisNone
	}
	
	// Convert mouse coordinates to normalized device coordinates (-1 to 1)
	ndcX := (2.0 * mouseX / float64(screenWidth)) - 1.0
	ndcY := 1.0 - (2.0 * mouseY / float64(screenHeight))
	
	// Create ray from camera through mouse position
	ray := g.createMouseRay(ndcX, ndcY, camera, view, projection)
	
	// Test ray intersection with gizmo axes
	selectedAxis := g.testRayIntersection(ray)
	g.selectedAxis = selectedAxis
	
	return selectedAxis
}

// Ray represents a 3D ray with origin and direction
type Ray struct {
	Origin    bmath.Vector3
	Direction bmath.Vector3
}

// createMouseRay creates a ray from the camera through the mouse position
func (g *Gizmo) createMouseRay(ndcX, ndcY float64, camera *bmath.Vector3, view, projection bmath.Matrix4) Ray {
	// Create inverse view-projection matrix
	viewProj := projection.Multiply(view)
	invViewProj := viewProj.Inverse()
	
	// Near and far points in normalized device coordinates
	nearPoint := bmath.NewVector4(float32(ndcX), float32(ndcY), -1.0, 1.0)
	farPoint := bmath.NewVector4(float32(ndcX), float32(ndcY), 1.0, 1.0)
	
	// Transform to world coordinates
	nearWorld := invViewProj.MultiplyVector4(nearPoint)
	farWorld := invViewProj.MultiplyVector4(farPoint)
	
	// Perspective divide
	if nearWorld.W != 0 {
		nearWorld.X /= nearWorld.W
		nearWorld.Y /= nearWorld.W
		nearWorld.Z /= nearWorld.W
	}
	if farWorld.W != 0 {
		farWorld.X /= farWorld.W
		farWorld.Y /= farWorld.W
		farWorld.Z /= farWorld.W
	}
	
	// Create ray
	origin := bmath.NewVector3(nearWorld.X, nearWorld.Y, nearWorld.Z)
	direction := bmath.NewVector3(farWorld.X - nearWorld.X, farWorld.Y - nearWorld.Y, farWorld.Z - nearWorld.Z)
	direction = direction.Normalize()
	
	return Ray{
		Origin:    origin,
		Direction: direction,
	}
}

// testRayIntersection tests which gizmo axis (if any) the ray intersects with
func (g *Gizmo) testRayIntersection(ray Ray) GizmoAxis {
	const threshold = 0.1 // Distance threshold for selection
	
	// Test X-axis (red arrow from origin to (1,0,0) * scale)
	xStart := g.position
	xEnd := bmath.NewVector3(g.position.X + g.scale, g.position.Y, g.position.Z)
	if g.distanceRayToLineSegment(ray, xStart, xEnd) < threshold {
		return GizmoAxisX
	}
	
	// Test Y-axis (green arrow from origin to (0,1,0) * scale)
	yStart := g.position
	yEnd := bmath.NewVector3(g.position.X, g.position.Y + g.scale, g.position.Z)
	if g.distanceRayToLineSegment(ray, yStart, yEnd) < threshold {
		return GizmoAxisY
	}
	
	// Test Z-axis (blue arrow from origin to (0,0,1) * scale)
	zStart := g.position
	zEnd := bmath.NewVector3(g.position.X, g.position.Y, g.position.Z + g.scale)
	if g.distanceRayToLineSegment(ray, zStart, zEnd) < threshold {
		return GizmoAxisZ
	}
	
	return GizmoAxisNone
}

// distanceRayToLineSegment calculates the minimum distance from a ray to a line segment
func (g *Gizmo) distanceRayToLineSegment(ray Ray, segStart, segEnd bmath.Vector3) float32 {
	// Vector from segment start to end
	segDir := bmath.NewVector3(segEnd.X - segStart.X, segEnd.Y - segStart.Y, segEnd.Z - segStart.Z)
	segLength := segDir.Length()
	if segLength < 1e-6 {
		// Degenerate segment, treat as point
		diff := bmath.NewVector3(segStart.X - ray.Origin.X, segStart.Y - ray.Origin.Y, segStart.Z - ray.Origin.Z)
		return diff.Length()
	}
	
	segDir = segDir.Normalize()
	
	// Vector from ray origin to segment start
	toSeg := bmath.NewVector3(segStart.X - ray.Origin.X, segStart.Y - ray.Origin.Y, segStart.Z - ray.Origin.Z)
	
	// Calculate closest points on both lines (if they were infinite)
	rayDotSeg := ray.Direction.Dot(segDir)
	rayDotToSeg := ray.Direction.Dot(toSeg)
	segDotToSeg := segDir.Dot(toSeg)
	
	denom := 1.0 - rayDotSeg * rayDotSeg
	var t1, t2 float32
	
	if denom < 1e-6 {
		// Lines are parallel
		t1 = 0
		t2 = segDotToSeg
	} else {
		t1 = (rayDotSeg * segDotToSeg - rayDotToSeg) / denom
		t2 = (segDotToSeg - rayDotSeg * rayDotToSeg) / denom
	}
	
	// Clamp t2 to segment bounds
	if t2 < 0 {
		t2 = 0
	} else if t2 > segLength {
		t2 = segLength
	}
	
	// Calculate closest points
	rayPoint := bmath.NewVector3(
		ray.Origin.X + ray.Direction.X * t1,
		ray.Origin.Y + ray.Direction.Y * t1,
		ray.Origin.Z + ray.Direction.Z * t1,
	)
	
	segPoint := bmath.NewVector3(
		segStart.X + segDir.X * t2,
		segStart.Y + segDir.Y * t2,
		segStart.Z + segDir.Z * t2,
	)
	
	// Distance between closest points
	diff := bmath.NewVector3(rayPoint.X - segPoint.X, rayPoint.Y - segPoint.Y, rayPoint.Z - segPoint.Z)
	return diff.Length()
}

func (g *Gizmo) Render(view, projection bmath.Matrix4) {
	if !g.visible {
		return
	}
	
	// Save current OpenGL state
	var currentDepthFunc int32
	gl.GetIntegerv(gl.DEPTH_FUNC, &currentDepthFunc)
	var currentLineWidth float32
	gl.GetFloatv(gl.LINE_WIDTH, &currentLineWidth)
	
	// Set gizmo rendering state
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)
	gl.LineWidth(3.0)
	// No need to change polygon mode - gizmos use gl.LINES primitives directly
	
	g.lineShader.Use()
	g.lineShader.SetMatrix4("view", &view[0])
	g.lineShader.SetMatrix4("projection", &projection[0])
	
	// Create model matrix with position and scale
	model := bmath.NewMatrix4Identity()
	model[12] = g.position.X
	model[13] = g.position.Y
	model[14] = g.position.Z
	model[0] = g.scale
	model[5] = g.scale
	model[10] = g.scale
	
	g.lineShader.SetMatrix4("model", &model[0])
	
	switch g.Type {
	case GizmoTranslate:
		g.renderTranslateGizmo()
	case GizmoRotate:
		// TODO: Implement rotation gizmo
	case GizmoScale:
		// TODO: Implement scale gizmo
	}
	
	// Restore OpenGL state
	gl.DepthFunc(uint32(currentDepthFunc))
	gl.LineWidth(currentLineWidth)
}

func (g *Gizmo) renderTranslateGizmo() {
	// Render X axis
	if g.hoveredAxis == GizmoAxisX || g.selectedAxis == GizmoAxisX {
		if g.selectedAxis == GizmoAxisX {
			g.lineShader.SetVec3("color", g.selectedColor.X, g.selectedColor.Y, g.selectedColor.Z)
		} else {
			g.lineShader.SetVec3("color", g.hoveredColor.X, g.hoveredColor.Y, g.hoveredColor.Z)
		}
	} else {
		g.lineShader.SetVec3("color", g.xColor.X, g.xColor.Y, g.xColor.Z)
	}
	g.xAxisMesh.Draw()
	
	// Render Y axis
	if g.hoveredAxis == GizmoAxisY || g.selectedAxis == GizmoAxisY {
		if g.selectedAxis == GizmoAxisY {
			g.lineShader.SetVec3("color", g.selectedColor.X, g.selectedColor.Y, g.selectedColor.Z)
		} else {
			g.lineShader.SetVec3("color", g.hoveredColor.X, g.hoveredColor.Y, g.hoveredColor.Z)
		}
	} else {
		g.lineShader.SetVec3("color", g.yColor.X, g.yColor.Y, g.yColor.Z)
	}
	g.yAxisMesh.Draw()
	
	// Render Z axis
	if g.hoveredAxis == GizmoAxisZ || g.selectedAxis == GizmoAxisZ {
		if g.selectedAxis == GizmoAxisZ {
			g.lineShader.SetVec3("color", g.selectedColor.X, g.selectedColor.Y, g.selectedColor.Z)
		} else {
			g.lineShader.SetVec3("color", g.hoveredColor.X, g.hoveredColor.Y, g.hoveredColor.Z)
		}
	} else {
		g.lineShader.SetVec3("color", g.zColor.X, g.zColor.Y, g.zColor.Z)
	}
	g.zAxisMesh.Draw()
	
	// Render plane indicators
	if g.Type == GizmoTranslate {
		// XY plane
		if g.hoveredAxis == GizmoAxisAll || g.selectedAxis == GizmoAxisAll {
			g.lineShader.SetVec3("color", g.hoveredColor.X, g.hoveredColor.Y, g.hoveredColor.Z)
		} else {
			g.lineShader.SetVec3("color", 1.0, 1.0, 0.0)
		}
		g.xyPlaneMesh.Draw()
		
		// XZ plane
		g.lineShader.SetVec3("color", 1.0, 0.0, 1.0)
		g.xzPlaneMesh.Draw()
		
		// YZ plane
		g.lineShader.SetVec3("color", 0.0, 1.0, 1.0)
		g.yzPlaneMesh.Draw()
	}
}

func (g *Gizmo) Cleanup() {
	if g.xAxisMesh != nil {
		g.xAxisMesh.Delete()
	}
	if g.yAxisMesh != nil {
		g.yAxisMesh.Delete()
	}
	if g.zAxisMesh != nil {
		g.zAxisMesh.Delete()
	}
	if g.xyPlaneMesh != nil {
		g.xyPlaneMesh.Delete()
	}
	if g.xzPlaneMesh != nil {
		g.xzPlaneMesh.Delete()
	}
	if g.yzPlaneMesh != nil {
		g.yzPlaneMesh.Delete()
	}
}

// CalculateGizmoRay calculates which axis (if any) the mouse ray intersects
func (g *Gizmo) CalculateHoveredAxis(rayOrigin, rayDir bmath.Vector3, viewportWidth, viewportHeight int) GizmoAxis {
	// This is a simplified version - in a real implementation you'd do proper
	// ray-cylinder intersection tests for the axes
	// For now, return none
	return GizmoAxisNone
}