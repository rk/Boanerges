package settings

// ReaderTheme values match legacy App\Enums\ReaderTheme.
var ReaderThemes = []string{
	"auto", "light", "autumn", "bumblebee", "cupcake", "emerald", "garden", "caramellatte",
	"lemonade", "nord", "pastel", "silk", "valentine", "winter", "dark", "coffee", "dracula",
	"forest", "luxury", "night", "sunset", "synthwave",
}

func isReaderTheme(theme string) bool {
	for _, t := range ReaderThemes {
		if t == theme {
			return true
		}
	}

	return false
}
