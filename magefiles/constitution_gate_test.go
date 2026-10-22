package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gray-lake-labs/scrum-master/internal/runtime"
	"github.com/gray-lake-labs/scrum-master/internal/scaffold"
	"github.com/gray-lake-labs/scrum-master/internal/spectool"
	"gopkg.in/yaml.v3"
)

// TestConstitutionProjectionMatchesThePinnedRelease holds
// internal/scaffold/parent-types.yaml to the pinned release's own
// design.yaml: each projected type's location and required fields equal the
// release's, so a pin move refreshes the projection or fails here (srd043,
// GH-186). It skips only when the release cache is absent; mage audit and
// the rigs run with it present.
func TestConstitutionProjectionMatchesThePinnedRelease(t *testing.T) {
	tag, err := spectool.Pinned("..")
	if err != nil {
		t.Fatal(err)
	}
	rel, err := runtime.Locate(tag)
	if err != nil {
		t.Fatal(err)
	}
	design, err := os.ReadFile(filepath.Join(rel.CoreRoot(), "docs", "constitutions", "design.yaml"))
	if os.IsNotExist(err) {
		t.Skipf("the release cache has no source at %s; run mage audit or tpm install --runtime first", rel.Src)
	}
	if err != nil {
		t.Fatal(err)
	}
	type declared struct {
		Location       string   `yaml:"location"`
		RequiredFields []string `yaml:"required_fields"`
	}
	var release struct {
		DocumentTypes map[string]declared `yaml:"document_types"`
	}
	if err := yaml.Unmarshal(design, &release); err != nil {
		t.Fatal(err)
	}
	projection, err := os.ReadFile(filepath.Join("..", "internal", "scaffold", "parent-types.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var projected struct {
		DocumentTypes map[string]declared `yaml:"document_types"`
	}
	if err := yaml.Unmarshal(projection, &projected); err != nil {
		t.Fatal(err)
	}
	for name, p := range projected.DocumentTypes {
		r, ok := release.DocumentTypes[name]
		if !ok {
			t.Errorf("the projection lists %s, which the release does not declare", name)
			continue
		}
		if p.Location != r.Location {
			t.Errorf("%s location: projection %q, release %q", name, p.Location, r.Location)
		}
		if len(p.RequiredFields) != len(r.RequiredFields) {
			t.Errorf("%s required fields: projection %v, release %v", name, p.RequiredFields, r.RequiredFields)
			continue
		}
		for i := range p.RequiredFields {
			if p.RequiredFields[i] != r.RequiredFields[i] {
				t.Errorf("%s required field %d: projection %q, release %q", name, i, p.RequiredFields[i], r.RequiredFields[i])
			}
		}
	}
	if _, err := scaffold.RenderConstitution(); err != nil {
		t.Fatal(err)
	}
}
