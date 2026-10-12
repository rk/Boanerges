<?php

namespace App\Services\Bible;

use App\Services\Bible\Import\BibleCanon;

class OsisBookId
{
    /** @var array<string, string> */
    private const ALIASES = [
        'exod' => 'exo',
        'deut' => 'deu',
        'josh' => 'jos',
        'judg' => 'jdg',
        'ruth' => 'rut',
        '1sam' => '1sa',
        '2sam' => '2sa',
        '1kgs' => '1ki',
        '2kgs' => '2ki',
        '1chr' => '1ch',
        '2chr' => '2ch',
        'esth' => 'est',
        'ezra' => 'ezr',
        'ps' => 'psa',
        'prov' => 'pro',
        'eccl' => 'ecc',
        'song' => 'sng',
        'songs' => 'sng',
        'songofsongs' => 'sng',
        'ezek' => 'ezk',
        'joel' => 'jol',
        'obad' => 'oba',
        'jonah' => 'jon',
        'nah' => 'nam',
        'zeph' => 'zep',
        'zech' => 'zec',
        'matt' => 'mat',
        'matthew' => 'mat',
        'mark' => 'mrk',
        'luke' => 'luk',
        'john' => 'jhn',
        'acts' => 'act',
        'phil' => 'php',
        '1cor' => '1co',
        '2cor' => '2co',
        '1thess' => '1th',
        '2thess' => '2th',
        '1tim' => '1ti',
        '2tim' => '2ti',
        'titus' => 'tit',
        'phlm' => 'phm',
        '1pet' => '1pe',
        '2pet' => '2pe',
        '1john' => '1jn',
        '2john' => '2jn',
        '3john' => '3jn',
        'jude' => 'jud',
        'revelation' => 'rev',
        'genesis' => 'gen',
        'exodus' => 'exo',
        'leviticus' => 'lev',
        'numbers' => 'num',
        'deuteronomy' => 'deu',
        'joshua' => 'jos',
        'judges' => 'jdg',
        '1samuel' => '1sa',
        '2samuel' => '2sa',
        '1kings' => '1ki',
        '2kings' => '2ki',
        '1chronicles' => '1ch',
        '2chronicles' => '2ch',
        'nehemiah' => 'neh',
        'esther' => 'est',
        'psalm' => 'psa',
        'psalms' => 'psa',
        'proverbs' => 'pro',
        'ecclesiastes' => 'ecc',
        'isaiah' => 'isa',
        'jeremiah' => 'jer',
        'lamentations' => 'lam',
        'ezekiel' => 'ezk',
        'daniel' => 'dan',
        'hosea' => 'hos',
        'amos' => 'amo',
        'obadiah' => 'oba',
        'micah' => 'mic',
        'nahum' => 'nam',
        'habakkuk' => 'hab',
        'zephaniah' => 'zep',
        'haggai' => 'hag',
        'zechariah' => 'zec',
        'malachi' => 'mal',
        'romans' => 'rom',
        '1corinthians' => '1co',
        '2corinthians' => '2co',
        'galatians' => 'gal',
        'ephesians' => 'eph',
        'philippians' => 'php',
        'colossians' => 'col',
        '1thessalonians' => '1th',
        '2thessalonians' => '2th',
        '1timothy' => '1ti',
        '2timothy' => '2ti',
        'philemon' => 'phm',
        'hebrews' => 'heb',
        'james' => 'jas',
        '1peter' => '1pe',
        '2peter' => '2pe',
    ];

    /** @var array<string, string> */
    private const DISPLAY_NAMES = [
        'gen' => 'Genesis',
        'exo' => 'Exodus',
        'lev' => 'Leviticus',
        'num' => 'Numbers',
        'deu' => 'Deuteronomy',
        'jos' => 'Joshua',
        'jdg' => 'Judges',
        'rut' => 'Ruth',
        '1sa' => '1 Samuel',
        '2sa' => '2 Samuel',
        '1ki' => '1 Kings',
        '2ki' => '2 Kings',
        '1ch' => '1 Chronicles',
        '2ch' => '2 Chronicles',
        'ezr' => 'Ezra',
        'neh' => 'Nehemiah',
        'est' => 'Esther',
        'job' => 'Job',
        'psa' => 'Psalms',
        'pro' => 'Proverbs',
        'ecc' => 'Ecclesiastes',
        'sng' => 'Song of Songs',
        'isa' => 'Isaiah',
        'jer' => 'Jeremiah',
        'lam' => 'Lamentations',
        'ezk' => 'Ezekiel',
        'dan' => 'Daniel',
        'hos' => 'Hosea',
        'jol' => 'Joel',
        'amo' => 'Amos',
        'oba' => 'Obadiah',
        'jon' => 'Jonah',
        'mic' => 'Micah',
        'nam' => 'Nahum',
        'hab' => 'Habakkuk',
        'zep' => 'Zephaniah',
        'hag' => 'Haggai',
        'zec' => 'Zechariah',
        'mal' => 'Malachi',
        'mat' => 'Matthew',
        'mrk' => 'Mark',
        'luk' => 'Luke',
        'jhn' => 'John',
        'act' => 'Acts',
        'rom' => 'Romans',
        '1co' => '1 Corinthians',
        '2co' => '2 Corinthians',
        'gal' => 'Galatians',
        'eph' => 'Ephesians',
        'php' => 'Philippians',
        'col' => 'Colossians',
        '1th' => '1 Thessalonians',
        '2th' => '2 Thessalonians',
        '1ti' => '1 Timothy',
        '2ti' => '2 Timothy',
        'tit' => 'Titus',
        'phm' => 'Philemon',
        'heb' => 'Hebrews',
        'jas' => 'James',
        '1pe' => '1 Peter',
        '2pe' => '2 Peter',
        '1jn' => '1 John',
        '2jn' => '2 John',
        '3jn' => '3 John',
        'jud' => 'Jude',
        'rev' => 'Revelation',
    ];

    /** @var array<string, list<string>>|null */
    private static ?array $lookupValues = null;

    public static function normalize(string $bookId): ?string
    {
        $key = self::normalizeKey($bookId);

        foreach (BibleCanon::books() as $book) {
            if ($book['osis'] === $key) {
                return $book['osis'];
            }
        }

        return self::ALIASES[$key] ?? null;
    }

    public static function isCanon(string $bookId): bool
    {
        return self::normalize($bookId) !== null;
    }

    public static function displayName(string $bookId): string
    {
        $canonical = self::normalize($bookId);

        if ($canonical === null) {
            return strtoupper($bookId);
        }

        return self::DISPLAY_NAMES[$canonical] ?? strtoupper($canonical);
    }

    /** @return array<string, string> */
    public static function displayNames(): array
    {
        return self::DISPLAY_NAMES;
    }

    /** @return list<string> */
    public static function lookupValues(string $bookId): array
    {
        $canonical = self::normalize($bookId);

        if ($canonical === null) {
            return [self::normalizeKey($bookId)];
        }

        self::$lookupValues ??= self::buildLookupValues();

        return self::$lookupValues[$canonical] ?? [$canonical];
    }

    private static function normalizeKey(string $bookId): string
    {
        return strtolower(str_replace([' ', '.', '-'], '', $bookId));
    }

    /** @return array<string, list<string>> */
    private static function buildLookupValues(): array
    {
        $values = [];

        foreach (BibleCanon::books() as $book) {
            $values[$book['osis']] = [$book['osis']];
        }

        foreach (self::ALIASES as $alias => $canonical) {
            $values[$canonical][] = $alias;
        }

        foreach ($values as $canonical => $aliases) {
            $values[$canonical] = array_values(array_unique($aliases));
        }

        return $values;
    }
}
