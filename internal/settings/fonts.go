package settings

// ReaderFontStack returns CSS font-family for bible/scribe reader text.
func ReaderFontStack(family string) string {
	if family == "sans-serif" {
		return `system-ui, -apple-system, "Segoe UI", sans-serif`
	}
	return `"Iowan Old Style", "Palatino Linotype", Palatino, Georgia, serif`
}
