package typescript

func stripRelativePrefixFromModulePath(modulePath []string) []string {
	i := 0
	for i < len(modulePath) && (modulePath[i] == "." || modulePath[i] == "..") {
		i++
	}
	return modulePath[i:]
}
