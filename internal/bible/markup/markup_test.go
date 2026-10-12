package markup

import "testing"

func TestGBFItalic(t *testing.T) {
	input := "Fear not, Abram, I <FI>am<Fi> a shield to thee, thy reward <FI>is<Fi> exceeding great."
	want := "Fear not, Abram, I <em>am</em> a shield to thee, thy reward <em>is</em> exceeding great."
	if got := DefaultFormatter().Format(input); got != want {
		t.Fatalf("got %q", got)
	}
}

func TestGBFBoldAndRed(t *testing.T) {
	got := DefaultFormatter().Format("<FB>bold<Fb> and <FR>red<Fr>")
	want := `<strong>bold</strong> and <span class="words-of-jesus">red</span>`
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestGBFSuperAndUnderline(t *testing.T) {
	got := DefaultFormatter().Format("word<FS>1<Fs> and <FU>underlined<Fu>")
	want := `word<sup>1</sup> and <span class="underline">underlined</span>`
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestOsisHi(t *testing.T) {
	got := DefaultFormatter().Format(`<hi type="italic">added</hi> and <hi type="bold">strong</hi>`)
	want := `<em>added</em> and <strong>strong</strong>`
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestOsisSemantic(t *testing.T) {
	input := `darkness <transChange type="added">was</transChange> upon the face of the deep. <divineName>God</divineName> moved`
	want := `darkness <em>was</em> upon the face of the deep. <span class="divine-name">God</span> moved`
	if got := DefaultFormatter().Format(input); got != want {
		t.Fatalf("got %q", got)
	}
}

func TestRemovesNotesAndStrongs(t *testing.T) {
	input := `In the beginning <w lemma="strong:H430">God</w><note>footnote</note> created`
	want := `In the beginning God created`
	if got := DefaultFormatter().Format(input); got != want {
		t.Fatalf("got %q", got)
	}
}

func TestEscapesHTML(t *testing.T) {
	got := DefaultFormatter().Format("Use <script> & <FI>am<Fi> here")
	want := `Use &lt;script&gt; &amp; <em>am</em> here`
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestPlainText(t *testing.T) {
	input := `<verse sID="Gen.1.1"/>In the beginning <w lemma="strong:H430">God</w> created<verse eID="Gen.1.1"/>`
	if got := DefaultFormatter().ToPlainText(input); got != "In the beginning God created" {
		t.Fatalf("got %q", got)
	}
}

func TestNestedGBF(t *testing.T) {
	got := DefaultFormatter().Format("<FI><FB>bold italic<Fb><Fi>")
	want := `<em><strong>bold italic</strong></em>`
	if got != want {
		t.Fatalf("got %q", got)
	}
}
