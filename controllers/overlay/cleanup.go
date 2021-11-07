package overlay

import (
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
)

func FindDanglingOverlays(e *v1.Environment) (paths []string) {
	var dangling []string
	state := e.Status.FileState
	for _, ec := range e.Spec.Components {
		for _, file := range state[ec.Name] {
			if file.Type != "overlay" {
				continue
			}
			exists := false
			for _, overlay := range ec.OverlayFiles {
				notMarkedForDeletion := !file.SoftDelete
				overlayName := common.ExtractNameFromPath(overlay.Path)
				isSameFile := file.Filename == overlayName
				if isSameFile && notMarkedForDeletion {
					exists = true
					break
				}
			}
			if !exists {
				path := fmt.Sprintf(
					"%s/%s",
					//il.EnvironmentComponentTerraformDirectoryPath(e.Spec.TeamName, e.Spec.EnvName, ec.Name),
					file.Filename,
				)
				dangling = append(dangling, path)
			}
		}
	}

	return dangling
}
