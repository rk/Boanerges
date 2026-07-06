<?php

namespace App\Services\Dictionary;

use Illuminate\Support\Facades\DB;

class Webster1828Dictionary
{
    /** @var array<string, string> */
    private const PART_OF_SPEECH_LABELS = [
        'n.' => 'noun',
        'v.' => 'verb',
        'a.' => 'adjective',
        'adv.' => 'adverb',
        'prep.' => 'preposition',
        'pron.' => 'pronoun',
        'conj.' => 'conjunction',
        'interj.' => 'interjection',
        'p.' => 'participle',
    ];

    public function __construct(
        private DictionaryService $dictionary,
    ) {}

    /**
     * @return array{
     *     found: bool,
     *     word: string,
     *     variants?: list<array{
     *         partOfSpeech: string|null,
     *         partOfSpeechExpanded: string,
     *         definitions: list<string>
     *     }>,
     *     message?: string
     * }
     */
    public function lookup(string $word): array
    {
        $wordKey = $this->normalizeWordKey($word);

        if ($wordKey === '') {
            return [
                'found' => false,
                'word' => $word,
                'message' => 'Enter a valid word to look up.',
            ];
        }

        $rows = DB::table('dictionary_entries')
            ->where('word_key', $wordKey)
            ->orderBy('id')
            ->get(['word', 'part_of_speech', 'definitions']);

        if ($rows->isEmpty()) {
            return [
                'found' => false,
                'word' => $word,
                'message' => "Word '{$word}' is not in this dictionary.",
            ];
        }

        return [
            'found' => true,
            'word' => (string) $rows->first()->word,
            'variants' => $rows->map(function ($row): array {
                $partOfSpeech = $row->part_of_speech !== null ? (string) $row->part_of_speech : null;

                return [
                    'partOfSpeech' => $partOfSpeech,
                    'partOfSpeechExpanded' => $this->expandPartOfSpeech($partOfSpeech),
                    'definitions' => json_decode((string) $row->definitions, true) ?? [],
                ];
            })->all(),
        ];
    }

    /**
     * @return list<string>
     */
    public function suggest(string $prefix, int $limit = 10): array
    {
        $normalized = $this->normalizeWordKey($prefix);

        if (strlen($normalized) < 2) {
            return [];
        }

        $rows = DB::table('dictionary_entries')
            ->select('word', 'word_key')
            ->where('word_key', 'like', $normalized . '%')
            ->orderBy('word_key')
            ->limit($limit * 3)
            ->get();

        $suggestions = [];
        $seen = [];

        foreach ($rows as $row) {
            $wordKey = (string) $row->word_key;

            if (isset($seen[$wordKey])) {
                continue;
            }

            $seen[$wordKey] = true;
            $suggestions[] = (string) $row->word;

            if (count($suggestions) >= $limit) {
                break;
            }
        }

        return $suggestions;
    }

    public function isReady(): bool
    {
        return $this->dictionary->isImported();
    }

    public function normalizeWordKey(string $word): string
    {
        $trimmed = trim($word);

        if ($trimmed === '') {
            return '';
        }

        $token = preg_split('/\s+/u', $trimmed)[0] ?? '';
        $token = preg_replace('/^[^\p{L}\p{N}\'-]+|[^\p{L}\p{N}\'-]+$/u', '', $token) ?? '';

        return mb_strtolower($token);
    }

    private function expandPartOfSpeech(?string $partOfSpeech): string
    {
        if ($partOfSpeech === null || $partOfSpeech === '') {
            return '';
        }

        return self::PART_OF_SPEECH_LABELS[$partOfSpeech] ?? $partOfSpeech;
    }
}
