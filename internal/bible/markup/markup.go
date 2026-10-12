package markup

import (
	"html"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const maxPasses = 4

const placeholderPrefix = "\x7Fboanerges-markup-"

type Converter interface {
	Convert(text string, placeholders map[string]string, formatter *Formatter) string
}

// Formatter converts verse markup (GBF/OSIS) to HTML and plain text.
type Formatter struct {
	converters []Converter
}

func DefaultFormatter() *Formatter {
	return &Formatter{converters: DefaultConverters()}
}

func NewFormatter(converters []Converter) *Formatter {
	return &Formatter{converters: converters}
}

func (f *Formatter) Format(text string) string {
	if text == "" {
		return ""
	}

	placeholders := map[string]string{}
	text = f.convert(text, placeholders)
	text = html.EscapeString(text)

	keys := placeholderKeys(placeholders)
	for i := len(keys) - 1; i >= 0; i-- {
		ph := keys[i]
		text = strings.ReplaceAll(text, html.EscapeString(ph), placeholders[ph])
	}

	return text
}

func (f *Formatter) ToPlainText(text string) string {
	if text == "" {
		return ""
	}

	placeholders := map[string]string{}
	plain := f.convert(text, placeholders)

	keys := placeholderKeys(placeholders)
	for i := len(keys) - 1; i >= 0; i-- {
		ph := keys[i]
		plain = strings.ReplaceAll(plain, ph, stripTags(placeholders[ph]))
	}

	plain = stripTags(plain)
	plain = html.UnescapeString(plain)
	plain = regexp.MustCompile(`\s+`).ReplaceAllString(plain, " ")

	return strings.TrimSpace(plain)
}

func (f *Formatter) convert(text string, placeholders map[string]string) string {
	for pass := 0; pass < maxPasses; pass++ {
		for _, converter := range f.converters {
			text = converter.Convert(text, placeholders, f)
		}
	}

	return text
}

func DefaultConverters() []Converter {
	return []Converter{
		removeConverter([]string{
			`<note\b[^>]*>.*?</note>`,
			`<figure\b[^>]*>.*?</figure>`,
			`<title\b[^>]*>.*?</title>`,
			`<rdgGroup\b[^>]*>.*?</rdgGroup>`,
			`<rdg\b[^>]*>.*?</rdg>`,
			`<catchWord\b[^>]*>.*?</catchWord>`,
			`<index\b[^>]*>.*?</index>`,
			`<milestone\b[^>]*/>`,
			`<milestone\b[^>]*>.*?</milestone>`,
			`<FN[^>]*>.*?<Fn>`,
			`<RF[^>]*>.*?<Rf>`,
			`<CM>.*?<Cm>`,
		}),
		&osisHiConverter{},
		&osisSemanticConverter{},
		pairTag("FI", "Fi", "em", ""),
		pairTag("FB", "Fb", "strong", ""),
		pairTag("FR", "Fr", "span", "words-of-jesus"),
		pairTag("FO", "Fo", "span", "ot-quote"),
		pairTag("FS", "Fs", "sup", ""),
		pairTag("FU", "Fu", "span", "underline"),
		unwrapConverter([]string{
			"p", "l", "lg", "q", "a", "w", "inscription", "mentioned", "name",
			"reference", "seg", "salute", "signed", "closer", "speech", "speaker",
			"list", "item", "table", "head", "row", "cell", "caption", "chapter",
			"div", "verse", "abbr", "hi", "transChange", "divineName", "foreign",
		}),
		removeConverter([]string{
			`<RX[^>]*>.*?</Rx>`,
			`<RP[^>]*>`,
			`<Rp[^>]*>`,
			`<H[^>]*>`,
			`<B[^>]*>`,
			`<D[^>]*>`,
			`<J[^>]*>`,
			`<P.>`,
			`<W[^>]*>`,
			`<S[^>]*>`,
			`<N[^>]*>`,
			`<C.>`,
			`<TS[^>]*>.*?</Ts>`,
		}),
	}
}

func placeholderKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return placeholderIndex(keys[i]) < placeholderIndex(keys[j])
	})

	return keys
}

func placeholderIndex(key string) int {
	key = strings.TrimPrefix(key, placeholderPrefix)
	key = strings.TrimSuffix(key, "\x7F")
	n, _ := strconv.Atoi(key)

	return n
}

func createPlaceholder(htmlFragment string, placeholders map[string]string) string {
	ph := placeholderPrefix + strconv.Itoa(len(placeholders)) + "\x7F"
	placeholders[ph] = htmlFragment

	return ph
}

func wrapContent(tag, content, class string, escape bool) string {
	inner := content
	if escape {
		inner = html.EscapeString(content)
	}
	if class == "" {
		return "<" + tag + ">" + inner + "</" + tag + ">"
	}

	return `<` + tag + ` class="` + html.EscapeString(class) + `">` + inner + `</` + tag + `>`
}

func stripTags(s string) string {
	re := regexp.MustCompile(`(?s)<[^>]*>`)

	return re.ReplaceAllString(s, "")
}

type removeConverter []string

func (c removeConverter) Convert(text string, _ map[string]string, _ *Formatter) string {
	for _, pattern := range c {
		re := regexp.MustCompile(`(?s)` + pattern)
		text = re.ReplaceAllString(text, "")
	}

	return text
}

type pairTagConverter struct {
	open, close, htmlTag, class string
}

func pairTag(open, close, htmlTag, class string) Converter {
	return pairTagConverter{open: open, close: close, htmlTag: htmlTag, class: class}
}

func (c pairTagConverter) Convert(text string, placeholders map[string]string, formatter *Formatter) string {
	pattern := `(?is)<` + regexp.QuoteMeta(c.open) + `>(.*?)<` + regexp.QuoteMeta(c.close) + `>`
	re := regexp.MustCompile(pattern)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		sub := re.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		inner := sub[1]
		if formatter != nil {
			inner = formatter.convert(inner, placeholders)
		}
		wrapped := wrapContent(c.htmlTag, inner, c.class, formatter == nil)

		return createPlaceholder(wrapped, placeholders)
	})
}

type osisHiConverter struct{}

func (osisHiConverter) Convert(text string, placeholders map[string]string, formatter *Formatter) string {
	re := regexp.MustCompile(`(?is)<hi\b[^>]*\btype=["']([^"']+)["'][^>]*>(.*?)</hi>`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		sub := re.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		typeMap := map[string][2]string{
			"italic":       {"em", ""},
			"bold":         {"strong", ""},
			"small-caps":   {"span", "small-caps"},
			"superscript":  {"sup", ""},
			"subscript":    {"sub", ""},
			"underline":    {"span", "underline"},
			"line-through": {"del", ""},
		}
		mapping, ok := typeMap[strings.ToLower(sub[1])]
		if !ok {
			return sub[2]
		}
		inner := sub[2]
		if formatter != nil {
			inner = formatter.convert(inner, placeholders)
		}
		wrapped := wrapContent(mapping[0], inner, mapping[1], formatter == nil)

		return createPlaceholder(wrapped, placeholders)
	})
}

type osisSemanticConverter struct{}

func (osisSemanticConverter) Convert(text string, placeholders map[string]string, formatter *Formatter) string {
	transRe := regexp.MustCompile(`(?is)<transChange\b[^>]*\btype=["']([^"']+)["'][^>]*>(.*?)</transChange>`)
	text = transRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := transRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		tag, class := "span", "trans-change"
		switch strings.ToLower(sub[1]) {
		case "added":
			tag, class = "em", ""
		case "removed":
			tag, class = "del", ""
		}
		inner := sub[2]
		if formatter != nil {
			inner = formatter.convert(inner, placeholders)
		}
		wrapped := wrapContent(tag, inner, class, formatter == nil)

		return createPlaceholder(wrapped, placeholders)
	})

	divineRe := regexp.MustCompile(`(?is)<divineName>(.*?)</divineName>`)
	text = divineRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := divineRe.FindStringSubmatch(match)
		inner := sub[1]
		if formatter != nil {
			inner = formatter.convert(inner, placeholders)
		}
		wrapped := wrapContent("span", inner, "divine-name", formatter == nil)

		return createPlaceholder(wrapped, placeholders)
	})

	foreignRe := regexp.MustCompile(`(?is)<foreign>(.*?)</foreign>`)
	return foreignRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := foreignRe.FindStringSubmatch(match)
		inner := sub[1]
		if formatter != nil {
			inner = formatter.convert(inner, placeholders)
		}
		wrapped := wrapContent("span", inner, "foreign", formatter == nil)

		return createPlaceholder(wrapped, placeholders)
	})
}

type unwrapConverter []string

func (c unwrapConverter) Convert(text string, _ map[string]string, _ *Formatter) string {
	for _, tag := range c {
		open := regexp.MustCompile(`(?i)<` + regexp.QuoteMeta(tag) + `\b[^>]*>`)
		close := regexp.MustCompile(`(?i)</` + regexp.QuoteMeta(tag) + `>`)
		text = open.ReplaceAllString(text, "")
		text = close.ReplaceAllString(text, "")
	}

	return text
}
