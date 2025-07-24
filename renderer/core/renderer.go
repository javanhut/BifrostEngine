package core

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/javanhut/BifrostEngine/m/v2/assets"
	"github.com/javanhut/BifrostEngine/m/v2/camera"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/opengl"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
)

type Renderer struct {
	window  *window.Window
	context *opengl.Context
	shader  *opengl.Shader
	materialShader *opengl.Shader
	lightingShader *opengl.Shader
	lineShader *opengl.Shader
	triangle *opengl.Mesh
	cube     *opengl.Mesh
	sphere   *opengl.Mesh
	cylinder *opengl.Mesh
	plane    *opengl.Mesh
	triangleMesh *opengl.Mesh
	pyramid  *opengl.Mesh
	// Lighting-enabled meshes
	cubeLighting *opengl.Mesh
	sphereLighting *opengl.Mesh
	planeLighting *opengl.Mesh
	triangleLighting *opengl.Mesh
	cylinderLighting *opengl.Mesh
	pyramidLighting *opengl.Mesh
	camera  *camera.Camera3D
	gridMesh *opengl.Mesh
	textureManager *TextureManager
	lightingSystem *LightingSystem
	assetManager *assets.AssetManager
	gizmo *Gizmo
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
	
	// Create material shader for texture support
	materialShader, err := opengl.NewShader(opengl.MaterialVertexShader, opengl.MaterialFragmentShader)
	if err != nil {
		win.Destroy()
		return nil, fmt.Errorf("failed to create material shader: %w", err)
	}
	
	// Create lighting shader for lighting effects
	lightingShader, err := opengl.NewShader(opengl.LightingVertexShader, opengl.LightingFragmentShader)
	if err != nil {
		win.Destroy()
		return nil, fmt.Errorf("failed to create lighting shader: %w", err)
	}
	
	// Create line shader for grid and gizmos - use existing shader for now
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
	
	// Create all mesh types with UV coordinates
	cube := opengl.NewCubeMeshWithUV()
	sphere := opengl.NewSphereMeshWithUV()
	cylinder := opengl.NewCylinderMeshWithUV()
	plane := opengl.NewPlaneMeshWithUV()
	triangleMesh := opengl.NewTriangleMeshWithUV()
	pyramid := opengl.NewPyramidMeshWithUV()
	
	// Create lighting-enabled meshes
	cubeLighting := opengl.NewCubeMeshWithLighting()
	sphereLighting := opengl.NewSphereMeshWithLighting()
	planeLighting := opengl.NewPlaneMeshWithLighting()
	triangleLighting := opengl.NewTriangleMeshWithLighting()
	cylinderLighting := opengl.NewCylinderMeshWithLighting()
	pyramidLighting := opengl.NewPyramidMeshWithLighting()

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
	
	// Create texture manager
	texManager := NewTextureManager()
	
	// Create lighting system
	lightingSystem := NewLightingSystem()
	
	// Create asset manager
	assetManager := assets.NewAssetManager()
	
	// Create gizmo
	gizmo := NewGizmo(lineShader)

	return &Renderer{
		window:  win,
		context: ctx,
		shader:  shader,
		materialShader: materialShader,
		lightingShader: lightingShader,
		lineShader: lineShader,
		triangle: triangle,
		cube:     cube,
		sphere:   sphere,
		cylinder: cylinder,
		plane:    plane,
		triangleMesh: triangleMesh,
		pyramid:  pyramid,
		cubeLighting: cubeLighting,
		sphereLighting: sphereLighting,
		planeLighting: planeLighting,
		triangleLighting: triangleLighting,
		cylinderLighting: cylinderLighting,
		pyramidLighting: pyramidLighting,
		camera:  cam,
		textureManager: texManager,
		lightingSystem: lightingSystem,
		assetManager: assetManager,
		gizmo: gizmo,
	}, nil
}

func (r *Renderer) BeginFrame() {
	r.BeginFrameWithMode(false) // Default to fill mode
}

func (r *Renderer) BeginFrameWithMode(wireframe bool) {
	width, height := r.window.GetSize()
	r.context.SetViewport(0, 0, int32(width), int32(height))
	r.context.Clear(0.1, 0.1, 0.1, 1.0)
	
	// Set proper OpenGL state for 3D rendering
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	
	// Set polygon mode based on wireframe setting
	if wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.LineWidth(1.5) // Slightly thicker lines for wireframe
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.LineWidth(1.0)
	}
	
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
	r.DrawCubeWithTextureToggle(model, false)
}

func (r *Renderer) DrawCubeWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.cube.Draw()
}

func (r *Renderer) DrawTriangleWithTransform(model bmath.Matrix4) {
	r.DrawTriangleWithTextureToggle(model, false)
}

func (r *Renderer) DrawTriangleWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.triangle.Draw()
}

func (r *Renderer) DrawSphereWithTransform(model bmath.Matrix4) {
	r.DrawSphereWithTextureToggle(model, false)
}

func (r *Renderer) DrawSphereWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.sphere.Draw()
}

func (r *Renderer) DrawCylinderWithTransform(model bmath.Matrix4) {
	r.DrawCylinderWithTextureToggle(model, false)
}

func (r *Renderer) DrawCylinderWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.cylinder.Draw()
}

func (r *Renderer) DrawPlaneWithTransform(model bmath.Matrix4) {
	r.DrawPlaneWithTextureToggle(model, false)
}

func (r *Renderer) DrawPlaneWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.plane.Draw()
}

func (r *Renderer) DrawTriangleMeshWithTransform(model bmath.Matrix4) {
	r.DrawTriangleMeshWithTextureToggle(model, false)
}

func (r *Renderer) DrawTriangleMeshWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.triangleMesh.Draw()
}

func (r *Renderer) DrawPyramidWithTransform(model bmath.Matrix4) {
	r.DrawPyramidWithTextureToggle(model, false)
}

func (r *Renderer) DrawPyramidWithTextureToggle(model bmath.Matrix4, useTextures bool) {
	r.materialShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.materialShader.SetMatrix4("model", &model[0])
	r.materialShader.SetMatrix4("view", &view[0])
	r.materialShader.SetMatrix4("projection", &projection[0])
	
	if useTextures {
		// Load and bind checkerboard texture
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.materialShader.SetBool("useTexture", true)
			r.materialShader.SetInt("diffuseTexture", 0)
		} else {
			r.materialShader.SetBool("useTexture", false)
		}
	} else {
		r.materialShader.SetBool("useTexture", false)
	}
	
	r.materialShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.materialShader.SetFloat("alpha", 1.0)
	
	r.pyramid.Draw()
}

// DrawCubeWithLighting renders a cube with full lighting effects
func (r *Renderer) DrawCubeWithLighting(model bmath.Matrix4, useTextures bool) {
	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix (inverse transpose of upper-left 3x3 of model matrix)
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties
	r.lightingShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 32.0)
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	r.cubeLighting.Draw()
}

// setLightingUniforms applies current lighting system settings to the lighting shader
func (r *Renderer) setLightingUniforms() {
	// Set ambient lighting
	ambient := r.lightingSystem.AmbientLight
	ambientColor := [3]float32{
		ambient.Color[0] * ambient.Intensity,
		ambient.Color[1] * ambient.Intensity,
		ambient.Color[2] * ambient.Intensity,
	}
	r.lightingShader.SetVec3("ambientLight", ambientColor[0], ambientColor[1], ambientColor[2])
	
	// Set directional light
	if len(r.lightingSystem.DirectionalLights) > 0 {
		dirLight := r.lightingSystem.DirectionalLights[0]
		r.lightingShader.SetBool("dirLightEnabled", dirLight.Enabled)
		if dirLight.Enabled {
			r.lightingShader.SetVec3("dirLightDirection", dirLight.Direction.X, dirLight.Direction.Y, dirLight.Direction.Z)
			r.lightingShader.SetVec3("dirLightColor", dirLight.Color[0], dirLight.Color[1], dirLight.Color[2])
			r.lightingShader.SetFloat("dirLightIntensity", dirLight.Intensity)
		}
	} else {
		r.lightingShader.SetBool("dirLightEnabled", false)
	}
	
	// Set point lights
	numPointLights := len(r.lightingSystem.PointLights)
	if numPointLights > 4 {
		numPointLights = 4 // Shader maximum
	}
	r.lightingShader.SetInt("numPointLights", int32(numPointLights))
	
	for i := 0; i < numPointLights; i++ {
		pointLight := r.lightingSystem.PointLights[i]
		if pointLight.Enabled {
			posUniform := fmt.Sprintf("pointLightPositions[%d]", i)
			colorUniform := fmt.Sprintf("pointLightColors[%d]", i)
			intensityUniform := fmt.Sprintf("pointLightIntensities[%d]", i)
			attenuationUniform := fmt.Sprintf("pointLightAttenuations[%d]", i)
			
			r.lightingShader.SetVec3(posUniform, pointLight.Position.X, pointLight.Position.Y, pointLight.Position.Z)
			r.lightingShader.SetVec3(colorUniform, pointLight.Color[0], pointLight.Color[1], pointLight.Color[2])
			r.lightingShader.SetFloat(intensityUniform, pointLight.Intensity)
			r.lightingShader.SetVec3(attenuationUniform, pointLight.Attenuation[0], pointLight.Attenuation[1], pointLight.Attenuation[2])
		}
	}
}

// DrawSphereWithLighting renders a sphere with full lighting effects
func (r *Renderer) DrawSphereWithLighting(model bmath.Matrix4, useTextures bool) {
	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties
	r.lightingShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 64.0) // Higher shininess for sphere
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	r.sphereLighting.Draw()
}

// DrawPlaneWithLighting renders a plane with full lighting effects
func (r *Renderer) DrawPlaneWithLighting(model bmath.Matrix4, useTextures bool) {
	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties
	r.lightingShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 16.0) // Lower shininess for plane
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	r.planeLighting.Draw()
}

// DrawTriangleWithLighting renders a triangle with full lighting effects
func (r *Renderer) DrawTriangleWithLighting(model bmath.Matrix4, useTextures bool) {
	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties
	r.lightingShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 8.0) // Low shininess for flat triangle
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	r.triangleLighting.Draw()
}

// DrawCylinderWithLighting renders a cylinder with full lighting effects
func (r *Renderer) DrawCylinderWithLighting(model bmath.Matrix4, useTextures bool) {
	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties
	r.lightingShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 32.0) // Medium shininess for cylinder
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	r.cylinderLighting.Draw()
}

// DrawPyramidWithLighting renders a pyramid with full lighting effects
func (r *Renderer) DrawPyramidWithLighting(model bmath.Matrix4, useTextures bool) {
	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties
	r.lightingShader.SetVec3("materialColor", 1.0, 1.0, 1.0)
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 16.0) // Medium-low shininess for pyramid
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	r.pyramidLighting.Draw()
}

// GetLightingSystem returns the lighting system for external configuration
func (r *Renderer) GetLightingSystem() *LightingSystem {
	return r.lightingSystem
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
	r.cubeLighting.Delete()
	r.sphereLighting.Delete()
	r.planeLighting.Delete()
	r.triangleLighting.Delete()
	r.cylinderLighting.Delete()
	r.pyramidLighting.Delete()
	r.shader.Delete()
	r.materialShader.Delete()
	r.lightingShader.Delete()
	r.lineShader.Delete()
	r.textureManager.Cleanup()
	r.assetManager.UnloadAll()
	if r.gridMesh != nil {
		r.gridMesh.Delete()
	}
	if r.gizmo != nil {
		r.gizmo.Cleanup()
	}
	r.window.Destroy()
}

// GetAssetManager returns the asset manager for loading meshes
func (r *Renderer) GetAssetManager() *assets.AssetManager {
	return r.assetManager
}

// DrawLoadedMesh renders a loaded OBJ mesh with lighting
func (r *Renderer) DrawLoadedMesh(filepath string, model bmath.Matrix4, useTextures bool) error {
	mesh, exists := r.assetManager.GetMesh(filepath)
	if !exists {
		return fmt.Errorf("mesh not loaded: %s", filepath)
	}

	r.lightingShader.Use()
	
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	// Set transformation matrices
	r.lightingShader.SetMatrix4("model", &model[0])
	r.lightingShader.SetMatrix4("view", &view[0])
	r.lightingShader.SetMatrix4("projection", &projection[0])
	
	// Calculate normal matrix
	normalMatrix := [9]float32{
		model[0], model[1], model[2],
		model[4], model[5], model[6],
		model[8], model[9], model[10],
	}
	r.lightingShader.SetMatrix3("normalMatrix", &normalMatrix[0])
	
	// Set camera position for specular calculations
	cameraPos := r.camera.GetPosition()
	r.lightingShader.SetVec3("viewPos", cameraPos.X, cameraPos.Y, cameraPos.Z)
	
	// Set material properties for loaded mesh
	r.lightingShader.SetVec3("materialColor", 0.9, 0.9, 0.9) // Slightly off-white
	r.lightingShader.SetFloat("alpha", 1.0)
	r.lightingShader.SetFloat("shininess", 32.0) // Higher shininess for imported models
	
	// Set texture properties
	if useTextures {
		textureID, err := r.textureManager.LoadCheckerboardTexture()
		if err == nil {
			r.textureManager.BindTexture(textureID, 0)
			r.lightingShader.SetBool("useTexture", true)
			r.lightingShader.SetInt("diffuseTexture", 0)
		} else {
			r.lightingShader.SetBool("useTexture", false)
		}
	} else {
		r.lightingShader.SetBool("useTexture", false)
	}
	
	// Set lighting uniforms
	r.setLightingUniforms()
	
	// Draw the loaded mesh
	mesh.Mesh.Draw()
	
	return nil
}

// LoadMesh loads a mesh file into the asset manager
func (r *Renderer) LoadMesh(filepath string) error {
	_, err := r.assetManager.LoadMesh(filepath)
	return err
}

// UnloadMesh unloads a mesh from the asset manager
func (r *Renderer) UnloadMesh(filepath string) {
	r.assetManager.UnloadMesh(filepath)
}

// GetLoadedMeshes returns a list of loaded mesh names
func (r *Renderer) GetLoadedMeshes() []string {
	return r.assetManager.GetLoadedMeshes()
}

// GetAssetStats returns statistics about loaded assets
func (r *Renderer) GetAssetStats() assets.AssetStats {
	return r.assetManager.GetStats()
}

// Gizmo-related methods
func (r *Renderer) RenderGizmo(position bmath.Vector3) {
	if r.gizmo == nil {
		return
	}
	
	r.gizmo.SetPosition(position)
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	r.gizmo.Render(view, projection)
}

func (r *Renderer) SetGizmoType(gizmoType GizmoType) {
	if r.gizmo != nil {
		r.gizmo.SetType(gizmoType)
	}
}

func (r *Renderer) SetGizmoVisible(visible bool) {
	if r.gizmo != nil {
		r.gizmo.SetVisible(visible)
	}
}

func (r *Renderer) SetGizmoScale(scale float32) {
	if r.gizmo != nil {
		r.gizmo.SetScale(scale)
	}
}

func (r *Renderer) GetGizmo() *Gizmo {
	return r.gizmo
}

// HandleGizmoMouseMove processes mouse movement for gizmo hover detection
func (r *Renderer) HandleGizmoMouseMove(mouseX, mouseY float64) {
	if r.gizmo == nil {
		return
	}
	
	width, height := r.window.GetSize()
	camera := r.camera.GetPosition()
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	r.gizmo.HandleMouseMove(mouseX, mouseY, &camera, view, projection, width, height)
}

// HandleGizmoMouseClick processes mouse clicks for gizmo selection
func (r *Renderer) HandleGizmoMouseClick(mouseX, mouseY float64) GizmoAxis {
	if r.gizmo == nil {
		return GizmoAxisNone
	}
	
	width, height := r.window.GetSize()
	camera := r.camera.GetPosition()
	view := r.camera.GetViewMatrix()
	projection := r.camera.GetProjectionMatrix()
	
	return r.gizmo.HandleMouseClick(mouseX, mouseY, &camera, view, projection, width, height)
}

// GetGizmoSelectedAxis returns the currently selected gizmo axis
func (r *Renderer) GetGizmoSelectedAxis() GizmoAxis {
	if r.gizmo == nil {
		return GizmoAxisNone
	}
	return r.gizmo.GetSelectedAxis()
}

// GetGizmoHoveredAxis returns the currently hovered gizmo axis
func (r *Renderer) GetGizmoHoveredAxis() GizmoAxis {
	if r.gizmo == nil {
		return GizmoAxisNone
	}
	return r.gizmo.GetHoveredAxis()
}

// GetCubeLightingMeshDebugInfo returns debug information about the cube lighting mesh
func (r *Renderer) GetCubeLightingMeshDebugInfo() (int32, uint32, bool) {
	if r.cubeLighting == nil {
		return 0, 0, false
	}
	return r.cubeLighting.IndexCount, r.cubeLighting.DrawMode, r.cubeLighting.IndexCount > 0
}