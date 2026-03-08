package adapter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/adapter"
)

func TestFileSkillInstaller_Install(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "skill-out")

	skill := []byte("# My Skill")
	schema := []byte(`{"type":"object"}`)

	installer := adapter.NewFileSkillInstaller()
	err := installer.Install(dir, skill, schema)
	require.NoError(t, err)

	gotSkill, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, skill, gotSkill)

	gotSchema, err := os.ReadFile(filepath.Join(dir, "walkthrough.schema.json"))
	require.NoError(t, err)
	assert.Equal(t, schema, gotSchema)
}

func TestFileSkillInstaller_Install_CreatesNestedDirs(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c")

	installer := adapter.NewFileSkillInstaller()
	err := installer.Install(dir, []byte("skill"), []byte("schema"))
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, "SKILL.md"))
	assert.NoError(t, err)
}

func TestFileSkillInstaller_Install_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("old"), 0644))

	installer := adapter.NewFileSkillInstaller()
	err := installer.Install(dir, []byte("new"), []byte("schema"))
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, []byte("new"), got)
}
