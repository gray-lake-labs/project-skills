package scaffold

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/gray-lake-labs/scrum-master/internal/spectool"
	"gopkg.in/yaml.v3"
)

// constitutionTemplate is the design constitution with a marker where each
// parent-declared type's location and required fields go; the articles and
// tpm's own types are authored prose. RenderConstitution fills the markers
// from parentTypes and the spec-tool overlay, so the constitution every
// program carries is produced from the declarations, never hand-synced
// (srd043, GH-186).
//
//go:embed constitution.md
var constitutionTemplate []byte

// parentTypes is the projection of the pinned parent release's design.yaml
// for the types the constitution lists: location and required fields. A gate
// test holds it equal to the release source in the cache.
//
//go:embed parent-types.yaml
var parentTypes []byte

// declaredType is what the renderer takes from a declaration.
type declaredType struct {
	Location       string   `yaml:"location"`
	RequiredFields []string `yaml:"required_fields"`
}

// declaredTypes reads a document_types projection: parent-types.yaml, or the
// spec-tool overlay, whose extra keys the projection shape ignores.
func declaredTypes(data []byte) (map[string]declaredType, error) {
	var f struct {
		DocumentTypes map[string]declaredType `yaml:"document_types"`
	}
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return f.DocumentTypes, nil
}

var generatedMarker = regexp.MustCompile(`^  <!-- generated: (location|required_fields) ([a-z_]+) -->$`)

// RenderConstitution produces docs/constitutions/design.md: the authored
// template with each parent-declared type's location and required fields
// filled from the pinned release's projection, the spec-tool overlay winning
// for the types it names (srd040 R1.6). It fails on a marker no declaration
// serves, so the constitution cannot silently drop a type.
func RenderConstitution() ([]byte, error) {
	types, err := declaredTypes(parentTypes)
	if err != nil {
		return nil, fmt.Errorf("parent-types.yaml: %w", err)
	}
	overlay, err := declaredTypes(spectool.Overlay)
	if err != nil {
		return nil, fmt.Errorf("spectool overlay: %w", err)
	}
	for name, d := range overlay {
		types[name] = d
	}
	template := string(constitutionTemplate)
	// The template opens with a note to its editors; the rendered
	// constitution starts at the frontmatter.
	if at := strings.Index(template, "-->\n"); strings.HasPrefix(template, "<!-- tpm:") && at >= 0 {
		template = template[at+len("-->\n"):]
	}
	var out []string
	for _, line := range strings.Split(template, "\n") {
		m := generatedMarker.FindStringSubmatch(line)
		if m == nil {
			out = append(out, line)
			continue
		}
		d, ok := types[m[2]]
		if !ok {
			return nil, fmt.Errorf("the constitution template names type %q, which no declaration serves", m[2])
		}
		switch m[1] {
		case "location":
			out = append(out, "  - location: "+d.Location)
		case "required_fields":
			if len(d.RequiredFields) == 0 {
				out = append(out, "  - required_fields: []")
				continue
			}
			out = append(out, "  - required_fields:")
			for _, f := range d.RequiredFields {
				out = append(out, "    1. "+f)
			}
		}
	}
	return []byte(strings.Join(out, "\n")), nil
}
