package pkg

import (
	"path"

	"github.com/ory/x/stringslice"
)

var supported = []string{"hydra", "kratos", "keto", "oathkeeper", "cli"}

func ProjectFromDir(dir string) string {
	project := path.Base(dir)
	if !stringslice.Has(supported, project) {
		Fatalf(`This command is expected to run in a directory named "hydra", "keto", "oathkeeper", "kratos".`)
		return ""
	}

	return project
}
