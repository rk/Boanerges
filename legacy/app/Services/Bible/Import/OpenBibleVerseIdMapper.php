<?php

namespace App\Services\Bible\Import;

use App\Services\Bible\OsisBookId;

class OpenBibleVerseIdMapper
{
    /** @var array<string, string> */
    private const OPENBIBLE_BOOKS = [
        'Gen' => 'gen',
        'Exod' => 'exo',
        'Lev' => 'lev',
        'Num' => 'num',
        'Deut' => 'deu',
        'Josh' => 'jos',
        'Judg' => 'jdg',
        'Ruth' => 'rut',
        '1Sam' => '1sa',
        '2Sam' => '2sa',
        '1Kgs' => '1ki',
        '2Kgs' => '2ki',
        '1Chr' => '1ch',
        '2Chr' => '2ch',
        'Ezra' => 'ezr',
        'Neh' => 'neh',
        'Esth' => 'est',
        'Job' => 'job',
        'Ps' => 'psa',
        'Prov' => 'pro',
        'Eccl' => 'ecc',
        'Song' => 'sng',
        'Isa' => 'isa',
        'Jer' => 'jer',
        'Lam' => 'lam',
        'Ezek' => 'ezk',
        'Dan' => 'dan',
        'Hos' => 'hos',
        'Joel' => 'jol',
        'Amos' => 'amo',
        'Obad' => 'oba',
        'Jonah' => 'jon',
        'Mic' => 'mic',
        'Nah' => 'nam',
        'Hab' => 'hab',
        'Zeph' => 'zep',
        'Hag' => 'hag',
        'Zech' => 'zec',
        'Mal' => 'mal',
        'Matt' => 'mat',
        'Mark' => 'mrk',
        'Luke' => 'luk',
        'John' => 'jhn',
        'Acts' => 'act',
        'Rom' => 'rom',
        '1Cor' => '1co',
        '2Cor' => '2co',
        'Gal' => 'gal',
        'Eph' => 'eph',
        'Phil' => 'php',
        'Col' => 'col',
        '1Thess' => '1th',
        '2Thess' => '2th',
        '1Tim' => '1ti',
        '2Tim' => '2ti',
        'Titus' => 'tit',
        'Phlm' => 'phm',
        'Heb' => 'heb',
        'Jas' => 'jas',
        '1Pet' => '1pe',
        '2Pet' => '2pe',
        '1John' => '1jn',
        '2John' => '2jn',
        '3John' => '3jn',
        'Jude' => 'jud',
        'Rev' => 'rev',
    ];

    /** @var array<int, array{book_id: string, chapter: int, verse: int}>|null */
    private static ?array $lookup = null;

    /** @var array<string, int>|null */
    private static ?array $osisVerseIndex = null;

    /**
     * @return array{book_id: string, chapter: int, verse: int}|null
     */
    public static function fromVerseId(int $verseId): ?array
    {
        self::$lookup ??= self::buildLookup();

        return self::$lookup[$verseId] ?? null;
    }

    /** @return array<int, array{book_id: string, chapter: int, verse: int}> */
    public static function all(): array
    {
        self::$lookup ??= self::buildLookup();

        return self::$lookup;
    }

    public static function normalizeOsisBookId(string $bookId): ?string
    {
        return OsisBookId::normalize($bookId);
    }

    /**
     * @return array{start: int, end: int|null}|null
     */
    public static function verseIdsFromReference(string $reference): ?array
    {
        $reference = trim($reference);

        if ($reference === '') {
            return null;
        }

        $parts = explode('-', $reference, 2);
        $startId = self::verseIdFromToken($parts[0]);

        if ($startId === null) {
            return null;
        }

        if (! isset($parts[1])) {
            return ['start' => $startId, 'end' => null];
        }

        $endId = self::verseIdFromToken($parts[1]);

        if ($endId === null) {
            return ['start' => $startId, 'end' => null];
        }

        return ['start' => $startId, 'end' => $endId];
    }

    /** @return array<int, array{book_id: string, chapter: int, verse: int}> */
    private static function buildLookup(): array
    {
        $lookup = [];
        $id = 1;

        foreach (BibleCanon::books() as $book) {
            foreach ($book['chapters'] as $chapterIndex => $verseCount) {
                $chapter = $chapterIndex + 1;

                for ($verse = 1; $verse <= $verseCount; $verse++) {
                    $lookup[$id] = [
                        'book_id' => $book['osis'],
                        'chapter' => $chapter,
                        'verse' => $verse,
                    ];
                    $id++;
                }
            }
        }

        return $lookup;
    }

    private static function verseIdFromToken(string $token): ?int
    {
        if (! preg_match('/^(\d?[A-Za-z]+)\.(\d+)\.(\d+)$/', trim($token), $matches)) {
            return null;
        }

        $osis = self::OPENBIBLE_BOOKS[$matches[1]] ?? null;

        if ($osis === null) {
            return null;
        }

        self::$osisVerseIndex ??= self::buildOsisVerseIndex();

        $key = $osis . '.' . $matches[2] . '.' . $matches[3];

        return self::$osisVerseIndex[$key] ?? null;
    }

    /** @return array<string, int> */
    private static function buildOsisVerseIndex(): array
    {
        $index = [];

        foreach (self::all() as $id => $ref) {
            $index["{$ref['book_id']}.{$ref['chapter']}.{$ref['verse']}"] = $id;
        }

        return $index;
    }
}
