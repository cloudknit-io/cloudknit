package terraformConfigModel

func buildModuleSource(moduleSource string) error {
	return nil
}

func buildModulePath(modulePath string) string {
	if modulePath == "" {
		return "."
	} else {
		return modulePath
	}
}
