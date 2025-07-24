package ui

import (
	"fmt"
	"time"
)

type Project struct {
	Name         string
	Path         string
	CreatedAt    time.Time
	LastModified time.Time
	SceneFile    string
}

type ProjectManager struct {
	projects       []Project
	currentProject *Project
	recentProjects []string
}

func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		projects:       make([]Project, 0),
		recentProjects: make([]string, 0),
	}
}

func (pm *ProjectManager) CreateProject(name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}
	
	project := Project{
		Name:         name,
		Path:         fmt.Sprintf("./projects/%s", name),
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
		SceneFile:    fmt.Sprintf("%s.bifrost", name),
	}
	
	pm.projects = append(pm.projects, project)
	pm.currentProject = &project
	pm.addToRecent(project.Name)
	
	return &project, nil
}

func (pm *ProjectManager) LoadProject(name string) (*Project, error) {
	for i := range pm.projects {
		if pm.projects[i].Name == name {
			pm.currentProject = &pm.projects[i]
			pm.addToRecent(name)
			return pm.currentProject, nil
		}
	}
	return nil, fmt.Errorf("project '%s' not found", name)
}

func (pm *ProjectManager) GetCurrentProject() *Project {
	return pm.currentProject
}

func (pm *ProjectManager) GetProjects() []Project {
	return pm.projects
}

func (pm *ProjectManager) GetRecentProjects() []string {
	return pm.recentProjects
}

func (pm *ProjectManager) addToRecent(name string) {
	// Remove if already exists
	for i, p := range pm.recentProjects {
		if p == name {
			pm.recentProjects = append(pm.recentProjects[:i], pm.recentProjects[i+1:]...)
			break
		}
	}
	
	// Add to front
	pm.recentProjects = append([]string{name}, pm.recentProjects...)
	
	// Keep only last 5
	if len(pm.recentProjects) > 5 {
		pm.recentProjects = pm.recentProjects[:5]
	}
}

func (pm *ProjectManager) SaveCurrentProject() error {
	if pm.currentProject == nil {
		return fmt.Errorf("no project is currently loaded")
	}
	
	pm.currentProject.LastModified = time.Now()
	// TODO: Implement actual file saving
	return nil
}