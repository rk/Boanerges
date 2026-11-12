package settings

const (
	ThemeAuto  = "auto"
	ThemeLight = "light"
	ThemeDark  = "dark"

	PoleMonastic = "monastic"
	PoleModern   = "modern"
)

var legacyLightThemes = map[string]struct{}{
	"light": {}, "sepia": {}, "autumn": {}, "bumblebee": {}, "cupcake": {}, "emerald": {},
	"garden": {}, "caramellatte": {}, "lemonade": {}, "nord": {}, "pastel": {}, "silk": {},
	"valentine": {}, "winter": {},
}

var legacyDarkThemes = map[string]struct{}{
	"dark": {}, "coffee": {}, "dracula": {}, "forest": {}, "luxury": {}, "night": {},
	"sunset": {}, "synthwave": {},
}

func isAppearanceTheme(theme string) bool {
	return theme == ThemeAuto || theme == ThemeLight || theme == ThemeDark
}

func isThemePole(pole string) bool {
	return pole == PoleMonastic || pole == PoleModern
}

func migrateLegacyTheme(theme string) string {
	if isAppearanceTheme(theme) {
		return theme
	}
	if _, ok := legacyLightThemes[theme]; ok {
		return ThemeLight
	}
	if _, ok := legacyDarkThemes[theme]; ok {
		return ThemeDark
	}
	return ThemeAuto
}

func normalizeThemePole(pole string) string {
	if isThemePole(pole) {
		return pole
	}
	return PoleMonastic
}
