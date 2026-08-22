package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

type Catalog struct {
	Root       string
	Objectives map[string]model.Objective
	Assertions map[string]model.Assertion
	Profiles   map[string]model.Profile
}

type objectiveDocument struct {
	SchemaVersion  string            `yaml:"schema_version"`
	CatalogVersion string            `yaml:"catalog_version"`
	Objectives     []model.Objective `yaml:"objectives"`
}

type assertionDocument struct {
	SchemaVersion  string            `yaml:"schema_version"`
	CatalogVersion string            `yaml:"catalog_version"`
	Assertions     []model.Assertion `yaml:"assertions"`
}

func Load(root string) (*Catalog, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve catalog root: %w", err)
	}
	c := &Catalog{
		Root:       abs,
		Objectives: map[string]model.Objective{},
		Assertions: map[string]model.Assertion{},
		Profiles:   map[string]model.Profile{},
	}
	if err := loadObjectives(c); err != nil {
		return nil, err
	}
	if err := loadAssertions(c); err != nil {
		return nil, err
	}
	if err := loadProfiles(c); err != nil {
		return nil, err
	}
	if err := c.validateReferences(); err != nil {
		return nil, err
	}
	return c, nil
}

func yamlFiles(root, directory string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(root, "catalog", directory, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("list %s catalog: %w", directory, err)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		return nil, fmt.Errorf("no %s catalog files under %s", directory, root)
	}
	return paths, nil
}

func decodeYAML(path string, destination any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func loadObjectives(c *Catalog) error {
	paths, err := yamlFiles(c.Root, "objectives")
	if err != nil {
		return err
	}
	for _, path := range paths {
		var document objectiveDocument
		if err := decodeYAML(path, &document); err != nil {
			return err
		}
		if document.SchemaVersion != "prc.objectives/v0.1" {
			return fmt.Errorf("unsupported objective schema %q in %s", document.SchemaVersion, path)
		}
		for _, objective := range document.Objectives {
			if objective.ID == "" || objective.Statement == "" || objective.Revision < 1 {
				return fmt.Errorf("invalid objective in %s", path)
			}
			if _, exists := c.Objectives[objective.ID]; exists {
				return fmt.Errorf("duplicate objective ID %s", objective.ID)
			}
			c.Objectives[objective.ID] = objective
		}
	}
	return nil
}

func loadAssertions(c *Catalog) error {
	paths, err := yamlFiles(c.Root, "assertions")
	if err != nil {
		return err
	}
	for _, path := range paths {
		var document assertionDocument
		if err := decodeYAML(path, &document); err != nil {
			return err
		}
		if document.SchemaVersion != "prc.assertions/v0.1" {
			return fmt.Errorf("unsupported assertion schema %q in %s", document.SchemaVersion, path)
		}
		for _, assertion := range document.Assertions {
			if assertion.ID == "" || assertion.ImplementationID == "" || assertion.Revision < 1 {
				return fmt.Errorf("invalid assertion in %s", path)
			}
			if _, exists := c.Assertions[assertion.ID]; exists {
				return fmt.Errorf("duplicate assertion ID %s", assertion.ID)
			}
			c.Assertions[assertion.ID] = assertion
		}
	}
	return nil
}

func loadProfiles(c *Catalog) error {
	paths, err := yamlFiles(c.Root, "profiles")
	if err != nil {
		return err
	}
	for _, path := range paths {
		var profile model.Profile
		if err := decodeYAML(path, &profile); err != nil {
			return err
		}
		if profile.SchemaVersion != "prc.profile/v0.1" || profile.ID == "" {
			return fmt.Errorf("invalid profile in %s", path)
		}
		if _, exists := c.Profiles[profile.ID]; exists {
			return fmt.Errorf("duplicate profile ID %s", profile.ID)
		}
		c.Profiles[profile.ID] = profile
	}
	return nil
}

func (c *Catalog) validateReferences() error {
	for _, objective := range c.Objectives {
		for _, assertionID := range objective.AssertionIDs {
			if _, ok := c.Assertions[assertionID]; !ok {
				return fmt.Errorf("objective %s references missing assertion %s", objective.ID, assertionID)
			}
		}
	}
	for _, assertion := range c.Assertions {
		for _, controlID := range assertion.ControlIDs {
			if _, ok := c.Objectives[controlID]; !ok {
				return fmt.Errorf("assertion %s references unavailable objective %s", assertion.ID, controlID)
			}
		}
	}
	for _, profile := range c.Profiles {
		for _, assertionID := range profile.AssertionIDs {
			if _, ok := c.Assertions[assertionID]; !ok {
				return fmt.Errorf("profile %s references missing assertion %s", profile.ID, assertionID)
			}
		}
	}
	return nil
}

func (c *Catalog) Profile(id string) (model.Profile, error) {
	profile, ok := c.Profiles[id]
	if !ok {
		return model.Profile{}, fmt.Errorf("unknown profile %q", id)
	}
	return profile, nil
}
