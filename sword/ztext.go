package sword

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"boanerges/internal/bible"
)

type blockMeta struct {
	compOffset  uint32
	compSize    uint32
	uncompSize  uint32
}

type verseEntry struct {
	blockNum uint16
	start    uint64
	size     uint16
}

type testamentData struct {
	bzv []byte
	bzs []byte
	bzz []byte

	verseIndex map[string]verseEntry // "Gen.1.1"
	blockCache map[uint16][]byte
	mu         sync.Mutex
}

// Module is an opened zText SWORD bible module (BOOK blocks + ZIP).
type Module struct {
	root string
	cfg  Config
	ot   *testamentData
	nt   *testamentData
}

// OpenModule opens a module from a .conf path (mods.d/asv.conf).
func OpenModule(confPath string) (*Module, error) {
	cfg, _, err := ParseConfFile(confPath)
	if err != nil {
		return nil, err
	}
	if cfg.ModDrv != "zText" {
		return nil, fmt.Errorf("unsupported ModDrv %q (need zText)", cfg.ModDrv)
	}
	if cfg.Compress != "ZIP" {
		return nil, fmt.Errorf("unsupported CompressType %q", cfg.Compress)
	}
	if cfg.BlockType != "BOOK" {
		return nil, fmt.Errorf("unsupported BlockType %q (need BOOK)", cfg.BlockType)
	}

	root := filepath.Clean(filepath.Join(filepath.Dir(confPath), "..", strings.TrimPrefix(cfg.DataPath, "./")))
	m := &Module{root: root, cfg: cfg}

	otDir := filepath.Join(root)
	if td, err := loadTestament(otDir, "ot"); err == nil && len(td.bzv) > 0 {
		m.ot = td
	}
	if td, err := loadTestament(otDir, "nt"); err == nil && len(td.bzv) > 0 {
		m.nt = td
	}

	if m.ot == nil && m.nt == nil {
		return nil, fmt.Errorf("no testament data under %s", root)
	}

	return m, nil
}

func loadTestament(dir, prefix string) (*testamentData, error) {
	read := func(name string) ([]byte, error) {
		return os.ReadFile(filepath.Join(dir, name))
	}
	bzv, err := read(prefix + ".bzv")
	if err != nil {
		return nil, err
	}
	bzs, err := read(prefix + ".bzs")
	if err != nil {
		return nil, err
	}
	bzz, err := read(prefix + ".bzz")
	if err != nil {
		return nil, err
	}

	td := &testamentData{
		bzv:        bzv,
		bzs:        bzs,
		bzz:        bzz,
		blockCache: map[uint16][]byte{},
	}
	idx, err := buildVerseIndex(bzv, prefix)
	if err != nil {
		return nil, err
	}
	td.verseIndex = idx

	return td, nil
}

func buildVerseIndex(bzv []byte, testament string) (map[string]verseEntry, error) {
	books := KVJBooks()
	startBook, endBook := 0, 39
	if testament == "nt" {
		startBook, endBook = 39, 66
	}

	off := 20 // module + testament intros (2 × 10 bytes per testament file)
	out := map[string]verseEntry{}

	for b := startBook; b < endBook; b++ {
		book := books[b]
		off += 10 // book heading
		for c := 0; c < book.MaxChapter; c++ {
			off += 10 // chapter heading
			for v := 0; v < book.Verses[c]; v++ {
				blockNum, err := readU16LE(bzv, off)
				if err != nil {
					return nil, err
				}
				start, err := readU48LE(bzv, off+2)
				if err != nil {
					return nil, err
				}
				size, err := readU16LE(bzv, off+8)
				if err != nil {
					return nil, err
				}
				off += 10

				if size > 0 {
					key := fmt.Sprintf("%s.%d.%d", book.CanonOSIS, c+1, v+1)
					out[key] = verseEntry{blockNum: blockNum, start: start, size: size}
				}
			}
		}
	}

	return out, nil
}

func parseBlocks(bzs []byte) ([]blockMeta, error) {
	var blocks []blockMeta
	for off := 0; off+12 <= len(bzs); off += 12 {
		compOff, err := readU32LE(bzs, off)
		if err != nil {
			return nil, err
		}
		compSize, err := readU32LE(bzs, off+4)
		if err != nil {
			return nil, err
		}
		uncompSize, err := readU32LE(bzs, off+8)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, blockMeta{compOffset: compOff, compSize: compSize, uncompSize: uncompSize})
	}

	return blocks, nil
}

func (td *testamentData) decompressBlock(blockNum uint16) ([]byte, error) {
	td.mu.Lock()
	if cached, ok := td.blockCache[blockNum]; ok {
		td.mu.Unlock()

		return cached, nil
	}
	td.mu.Unlock()

	blocks, err := parseBlocks(td.bzs)
	if err != nil {
		return nil, err
	}
	if int(blockNum) >= len(blocks) {
		return nil, fmt.Errorf("block %d out of range", blockNum)
	}
 blk := blocks[blockNum]
	end := int(blk.compOffset + blk.compSize)
	if end > len(td.bzz) {
		return nil, fmt.Errorf("compressed slice out of range")
	}
	compressed := td.bzz[blk.compOffset:end]

	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	td.mu.Lock()
	td.blockCache[blockNum] = decompressed
	td.mu.Unlock()

	return decompressed, nil
}

// GetVerse returns raw OSIS/XML markup for a reference (SWORD book abbrev, 1-based chapter/verse).
func (m *Module) GetVerse(bookAbbrev string, chapter, verse int) (string, error) {
	canon, ok := bible.NormalizeBookID(bookAbbrev)
	if !ok {
		canon = strings.ToLower(bookAbbrev)
	}
	key := fmt.Sprintf("%s.%d.%d", canon, chapter, verse)
	td := m.ot
	if isNTCanon(canon) {
		td = m.nt
	}
	if td == nil {
		return "", fmt.Errorf("testament data missing for %s", bookAbbrev)
	}

	entry, ok := td.verseIndex[key]
	if !ok || entry.size == 0 {
		return "", fmt.Errorf("verse not found: %s", key)
	}

	block, err := td.decompressBlock(entry.blockNum)
	if err != nil {
		return "", err
	}
	end := int(entry.start) + int(entry.size)
	if end > len(block) {
		return "", fmt.Errorf("verse slice out of range for %s", key)
	}

	return string(block[entry.start:end]), nil
}

func isNTCanon(canon string) bool {
	books := KVJBooks()
	for i, b := range books {
		if b.CanonOSIS == canon {
			return i >= 39
		}
	}

	return false
}

// Books lists module books.
func (m *Module) Books() []BookInfo {
	return KVJBooks()
}

// Config returns module configuration.
func (m *Module) Config() Config { return m.cfg }
