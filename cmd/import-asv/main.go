// Command import-asv builds extras/bundled/asv.sqlite from the bundled SWORD module.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
	"boanerges/sword"
)

func main() {
	repoRoot := flag.String("root", ".", "repository root")
	out := flag.String("out", "extras/bundled/asv.sqlite", "output sqlite path (translation tables only)")
	flag.Parse()

	conf := filepath.Join(*repoRoot, "extras/sword/mods.d/asv.conf")
	mod, err := sword.OpenModule(conf)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		log.Fatal(err)
	}
	_ = os.Remove(*out)

	mgr, err := schema.Open(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer mgr.Close()

	if err := mgr.CreateTables("asv"); err != nil {
		log.Fatal(err)
	}

	db := mgr.DB()
	formatter := markup.DefaultFormatter()
	books := mod.Books()

	for i, book := range books {
		testament := "ot"
		if i >= 39 {
			testament = "nt"
		}

		res, err := db.Exec(
			`INSERT INTO asv_books (name, osis_id, testament, chapters) VALUES (?, ?, ?, ?)`,
			book.Name, book.CanonOSIS, testament, book.MaxChapter,
		)
		if err != nil {
			log.Fatal(err)
		}
		bookID, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		for ch := 1; ch <= book.MaxChapter; ch++ {
			for v := 1; v <= book.Verses[ch-1]; v++ {
				text, err := mod.GetVerse(book.CanonOSIS, ch, v)
				if err != nil {
					log.Fatalf("%s %d:%d: %v", book.CanonOSIS, ch, v, err)
				}
				plain := formatter.ToPlainText(text)
				if _, err := db.Exec(
					`INSERT INTO asv_verses (book_id, chapter, verse, text, plain_text) VALUES (?, ?, ?, ?, ?)`,
					bookID, ch, v, text, plain,
				); err != nil {
					log.Fatal(err)
				}
			}
		}

		fmt.Printf("imported %s\n", book.CanonOSIS)
	}

	if err := mgr.CreateFtsTable("asv"); err != nil {
		log.Fatal(err)
	}
	if err := mgr.RebuildFtsIndex("asv", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote", *out)
}
