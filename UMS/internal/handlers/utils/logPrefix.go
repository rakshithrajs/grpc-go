package handlers

// LogPrefix returns a formatted log prefix string for the given function name.
func LogPrefix(fnName string) string {
	return "[" + fnName + "]: "
}
