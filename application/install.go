package application

import (
	"path/filepath"

	"github.com/4thel00z/code-walkthrough/port"
)

type InstallSkillUseCase struct {
	installer port.SkillInstaller
	skill     []byte
	schema    []byte
}

func NewInstallSkillUseCase(installer port.SkillInstaller, skill, schema []byte) *InstallSkillUseCase {
	return &InstallSkillUseCase{
		installer: installer,
		skill:     skill,
		schema:    schema,
	}
}

func (uc *InstallSkillUseCase) Install(dir string) error {
	return uc.installer.Install(dir, uc.skill, uc.schema)
}

func (uc *InstallSkillUseCase) DefaultInstallDir() string {
	return filepath.Join(".claude", "skills", "code-walkthrough")
}
