// Package scrum-master holds the Agent Plugins 1.0 package of the skills the
// tpm binary embeds; tpm install writes them for the developer's agent
// (srd037).
package scrummaster

import "embed"

// Plugin is the package plugin/: plugin.json and skills/<name>/SKILL.md.
//
//go:embed plugin
var Plugin embed.FS

// SkillsDir is where the skills sit in Plugin.
const SkillsDir = "plugin/skills"

// Pins is the parent release this corpus's format follows; the corpora tpm
// scaffolds pin the same release (srd043 R3.3).
//
//go:embed docs/external-pins.yaml
var Pins []byte

// Charter is the corpus charter this corpus gates on; the corpora tpm
// scaffolds gate on the same checks (srd043 R1.1).
//
//go:embed docs/corpus-charter.yaml
var Charter []byte
