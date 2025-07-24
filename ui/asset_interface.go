package ui

// AssetRenderer interface defines the methods needed for asset management in the GUI
type AssetRenderer interface {
	GetAssetStats() AssetStats
	GetAssetManager() AssetManager
	GetLoadedMeshes() []string
	LoadMesh(filepath string) error
	UnloadMesh(filepath string)
}

// AssetStats represents statistics about loaded assets
type AssetStats struct {
	LoadedMeshes  int
	TotalVertices int
	TotalIndices  int
}

// AssetManager interface for managing assets
type AssetManager interface {
	UnloadAll()
	LoadMesh(filepath string) (interface{}, error)
	GetLoadedMeshes() []string
	GetStats() AssetStats
}