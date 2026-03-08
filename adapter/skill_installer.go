package adapter

import (
	"os"
	"path/filepath"
)

type FileSkillInstaller struct{}

func NewFileSkillInstaller() *FileSkillInstaller {
	return &FileSkillInstaller{}
}

func (f *FileSkillInstaller) Install(dir string, skill []byte, schema []byte) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), skill, 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "walkthrough.schema.json"), schema, 0644)
}
