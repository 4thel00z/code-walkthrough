package skilldata

import _ "embed"

//go:embed skill.md
var SkillMD []byte

//go:embed walkthrough.schema.json
var SchemaJSON []byte
