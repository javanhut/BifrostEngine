package assets

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// AssetManager manages loading and caching of assets
type AssetManager struct {
	objLoader    *OBJLoader
	loadedMeshes map[string]*OBJMesh // Cache of loaded meshes
	mutex        sync.RWMutex       // Thread-safe access
}

// NewAssetManager creates a new asset manager
func NewAssetManager() *AssetManager {
	return &AssetManager{
		objLoader:    NewOBJLoader(),
		loadedMeshes: make(map[string]*OBJMesh),
	}
}

// LoadMesh loads a mesh from file, returns cached version if already loaded
func (am *AssetManager) LoadMesh(filePath string) (*OBJMesh, error) {
	am.mutex.RLock()
	if mesh, exists := am.loadedMeshes[filePath]; exists {
		am.mutex.RUnlock()
		return mesh, nil
	}
	am.mutex.RUnlock()

	// Determine file type and load
	ext := strings.ToLower(filepath.Ext(filePath))
	var mesh *OBJMesh
	var err error

	switch ext {
	case ".obj":
		mesh, err = am.objLoader.LoadOBJ(filePath)
	default:
		return nil, fmt.Errorf("unsupported mesh format: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	// Cache the loaded mesh
	am.mutex.Lock()
	am.loadedMeshes[filePath] = mesh
	am.mutex.Unlock()

	return mesh, nil
}

// GetLoadedMeshes returns a list of all loaded mesh names
func (am *AssetManager) GetLoadedMeshes() []string {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	var names []string
	for path, mesh := range am.loadedMeshes {
		names = append(names, fmt.Sprintf("%s (%s)", mesh.Name, path))
	}
	return names
}

// UnloadMesh removes a mesh from cache and cleans up resources
func (am *AssetManager) UnloadMesh(filePath string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if mesh, exists := am.loadedMeshes[filePath]; exists {
		mesh.Cleanup()
		delete(am.loadedMeshes, filePath)
	}
}

// UnloadAll cleans up all loaded assets
func (am *AssetManager) UnloadAll() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	for _, mesh := range am.loadedMeshes {
		mesh.Cleanup()
	}
	am.loadedMeshes = make(map[string]*OBJMesh)
}

// GetMesh returns a cached mesh if it exists
func (am *AssetManager) GetMesh(filePath string) (*OBJMesh, bool) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	mesh, exists := am.loadedMeshes[filePath]
	return mesh, exists
}

// SetOBJLoaderConfig configures the OBJ loader
func (am *AssetManager) SetOBJLoaderConfig(defaultColor [3]float32, flipUV, generateNormals bool) {
	am.objLoader.DefaultColor = defaultColor
	am.objLoader.FlipUV = flipUV
	am.objLoader.GenerateNormals = generateNormals
}

// GetStats returns statistics about loaded assets
func (am *AssetManager) GetStats() AssetStats {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	stats := AssetStats{
		LoadedMeshes: len(am.loadedMeshes),
	}

	// Calculate total vertices and indices
	for _, mesh := range am.loadedMeshes {
		stats.TotalVertices += len(mesh.Vertices) / 11 // 11 floats per vertex
		stats.TotalIndices += len(mesh.Indices)
	}

	return stats
}

// AssetStats contains statistics about loaded assets
type AssetStats struct {
	LoadedMeshes   int
	TotalVertices  int
	TotalIndices   int
}