package application_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/application"
)

type mockSkillInstaller struct {
	installedDir    string
	installedSkill  []byte
	installedSchema []byte
	err             error
}

func (m *mockSkillInstaller) Install(dir string, skill []byte, schema []byte) error {
	m.installedDir = dir
	m.installedSkill = skill
	m.installedSchema = schema
	return m.err
}

func TestInstallSkillUseCase_Install(t *testing.T) {
	skill := []byte("# Skill content")
	schema := []byte(`{"type":"object"}`)
	installer := &mockSkillInstaller{}

	uc := application.NewInstallSkillUseCase(installer, skill, schema)

	err := uc.Install("/tmp/test-dir")
	require.NoError(t, err)

	assert.Equal(t, "/tmp/test-dir", installer.installedDir)
	assert.Equal(t, skill, installer.installedSkill)
	assert.Equal(t, schema, installer.installedSchema)
}

func TestInstallSkillUseCase_Install_Error(t *testing.T) {
	installer := &mockSkillInstaller{err: errors.New("write failed")}

	uc := application.NewInstallSkillUseCase(installer, []byte("x"), []byte("y"))

	err := uc.Install("/tmp/fail")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "write failed")
}

func TestInstallSkillUseCase_DefaultInstallDir(t *testing.T) {
	uc := application.NewInstallSkillUseCase(nil, nil, nil)
	dir := uc.DefaultInstallDir()

	expected := filepath.Join(".claude", "skills", "code-walkthrough")
	assert.Equal(t, expected, dir)
}
