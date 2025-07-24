package core

import (
	"fmt"

	"github.com/javanhut/BifrostEngine/m/v2/camera"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/opengl"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
)

type Renderer struct {
	window  *window.Window
	context *opengl.Context
	shader  *opengl.Shader
	lineShader *opengl.Shader
	triangle *opengl.Mesh
	cube     *opengl.Mesh
	sphere   *opengl.Mesh
	cylinder *opengl.Mesh
	plane    *opengl.Mesh
	triangleMesh *opengl.Mesh
	pyramid  *opengl.Mesh
	camera  *camera.Camera3D
	gridMesh *opengl.Mesh
}

func New(width, height int, title string) (*Renderer, error) {
	win, err := window.New(width, height, title)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	ctx, err := opengl.NewContext()
	if err != nil {
		win.Destroy()
		return nil, fmt.Errorf("failed to create OpenGL context: %w", err)
	}

	shader, err := opengl.NewShader(opengl.DefaultVertexShader, opengl.DefaultFragmentShader)
	if err != nil {
		win.Destroy()
		return nil, fmt.Errorf("failed to create shader: %w", err)
	}
	
	// Create line shader for grid
	lineShader, err := opengl.NewShader(opengl.DefaultVertexShader, opengl.DefaultFragmentShader)
	if err != nil {
		win.Destroy()
		return nil, fmt.Errorf("failed to create line shader: %w", err)
	}

	// Triangle vertices: x, y, z, r, g, b
	// Simple setup: triangle at origin
	triangleVertices := []float32{
		// Position       Color
		-0.5, -0.5, 0.0,  1.0, 0.0, 0.0,  // Bottom left - Red
		 0.5, -0.5, 0.0,  0.0, 1.0, 0.0,  // Bottom right - Green
		 0.0,  0.5, 0.0,  0.0, 0.0, 1.0,  // Top - Blue
	}
	triangle := opengl.NewMesh(triangleVertices)
	
	// Create all mesh types
	cube := opengl.NewCubeMesh()
	sphere := opengl.NewSphereMesh()
	cylinder := opengl.NewCylinderMesh()
	plane := opengl.NewPlaneMesh()
	triangleMesh := opengl.NewTriangleMesh()
	pyramid := opengl.NewPyramidMesh()

	// Create camera back at z=3 looking at origin (standard setup)
	cameraPos := bmath.NewVector3(0, 0, 3)
	cameraTarget := bmath.NewVector3(0, 0, 0)
	cam := camera.NewCamera3D(
		cameraPos,
		cameraTarget,
		bmath.Radians(45),
		float32(width)/float32(height),
		1.0,  // Increase near plane
		100.0,
	)

	return &Renderer{
		window:  win,
		context: ctx,
		shader:  shader,
		lineShader: lineShader,
		triangle: triangle,
		cube:     cube,
		sphere:   sphere,
		cylinder: cylinder,
		plane:    plane,
		triangleMesh: triangleMesh,
		pyramid:  pyramid,
		camera:  cam,
	}, nil
}

func (r *Renderer) BeginFrame() {
	width, height := r.window.GetSize()
	r.context.SetViewport(0, 0, int32(width), int32(height))
	r.context.Clear(0.1, 0.1, 0.1, 1.0)
	
	// Update camera aspect ratio if window resized
	if width > 0 && height > 0 {
		r.camera.SetAspect(float32(width) / float32(height))
	}
}

func (r *Renderer) DrawTriangle() {
	r.shader.Use()
	
	// Set matrices
	model := bmath.NewMatrix4Identity()
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.triangle.Draw()
}

func (r *Renderer) DrawRotatingTriangle(time float32) {
	r.shader.Use()
	
	// Create rotation matrices for triangle
	rotZ := bmath.NewRotationZ(time)
	model := rotZ
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.triangle.Draw()
}

func (r *Renderer) DrawCube() {
	r.shader.Use()
	
	// Set matrices
	model := bmath.NewMatrix4Identity()
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.cube.Draw()
}

func (r *Renderer) DrawRotatingCube(time float32) {
	r.shader.Use()
	
	// Create rotation matrices
	rotY := bmath.NewRotationY(time)
	rotX := bmath.NewRotationX(time * 0.5)
	model := rotY.Multiply(rotX)
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.cube.Draw()
}

func (r *Renderer) DrawCubeWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.cube.Draw()
}

func (r *Renderer) DrawTriangleWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.triangle.Draw()
}

func (r *Renderer) DrawSphereWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.sphere.Draw()
}

func (r *Renderer) DrawCylinderWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.cylinder.Draw()
}

func (r *Renderer) DrawPlaneWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.plane.Draw()
}

func (r *Renderer) DrawTriangleMeshWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.triangleMesh.Draw()
}

func (r *Renderer) DrawPyramidWithTransform(model bmath.Matrix4) {
	r.shader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.shader.SetMatrix4("model", &model[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &projection[0])
	
	r.pyramid.Draw()
}

func (r *Renderer) GetCamera() *camera.Camera3D {
	return r.camera
}

func (r *Renderer) GetWindow() *window.Window {
	return r.window
}

func (r *Renderer) DrawTriangleNoCamera() {
	r.shader.Use()
	
	// Set identity matrices - no camera transformation
	identity := bmath.NewMatrix4Identity()
	
	r.shader.SetMatrix4("model", &identity[0])
	r.shader.SetMatrix4("view", &identity[0])
	r.shader.SetMatrix4("projection", &identity[0])
	
	r.triangle.Draw()
}

func (r *Renderer) DrawTriangleProjectionOnly() {
	r.shader.Use()
	
	// Use orthographic projection for mode 2 instead of perspective
	// This makes it easier to see the effect of projection alone
	ortho := bmath.NewOrthographic(-1, 1, -1, 1, -10, 10)
	identity := bmath.NewMatrix4Identity()
	
	r.shader.SetMatrix4("model", &identity[0])
	r.shader.SetMatrix4("view", &identity[0])
	r.shader.SetMatrix4("projection", &ortho[0])
	
	r.triangle.Draw()
}

func (r *Renderer) DrawTriangleViewOnly() {
	r.shader.Use()
	
	identity := bmath.NewMatrix4Identity()
	view := r.camera.GetViewMatrix()
	// Use orthographic projection for view-only mode
	ortho := bmath.NewOrthographic(-2, 2, -2, 2, -10, 10)
	
	r.shader.SetMatrix4("model", &identity[0])
	r.shader.SetMatrix4("view", &view[0])
	r.shader.SetMatrix4("projection", &ortho[0])
	
	r.triangle.Draw()
}

func (r *Renderer) DrawGrid(lines []bmath.Vector3, color [4]float32) {
	if len(lines) < 2 {
		return
	}
	
	// Convert lines to vertex data with color
	vertices := make([]float32, 0, len(lines)*6)
	for _, line := range lines {
		vertices = append(vertices, line.X, line.Y, line.Z, color[0], color[1], color[2])
	}
	
	// Create temporary mesh for grid lines
	if r.gridMesh != nil {
		r.gridMesh.Delete()
	}
	r.gridMesh = opengl.NewMeshLines(vertices)
	
	r.lineShader.Use()
	
	model := bmath.NewMatrix4Identity()
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.lineShader.SetMatrix4("model", &model[0])
	r.lineShader.SetMatrix4("view", &view[0])
	r.lineShader.SetMatrix4("projection", &projection[0])
	
	r.gridMesh.DrawLines()
}

func (r *Renderer) EndFrame() {
	r.window.SwapBuffers()
	r.window.PollEvents()
}

func (r *Renderer) ShouldClose() bool {
	return r.window.ShouldClose()
}

func (r *Renderer) Cleanup() {
	r.triangle.Delete()
	r.cube.Delete()
	r.sphere.Delete()
	r.cylinder.Delete()
	r.plane.Delete()
	r.triangleMesh.Delete()
	r.pyramid.Delete()
	r.shader.Delete()
	r.lineShader.Delete()
	if r.gridMesh != nil {
		r.gridMesh.Delete()
	}
	r.window.Destroy()
}