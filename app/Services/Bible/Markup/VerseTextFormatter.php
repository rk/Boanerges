<?php

namespace App\Services\Bible\Markup;

final class VerseTextFormatter
{
    private const int MAX_PASSES = 4;

    /**
     * @param  list<VerseMarkupConverter>  $converters
     */
    public function __construct(
        private array $converters,
    ) {}

    public function format(string $text): string
    {
        if ($text === '') {
            return '';
        }

        /** @var array<string, string> $placeholders */
        $placeholders = [];

        $text = $this->convert($text, $placeholders);

        $text = htmlspecialchars($text, ENT_QUOTES | ENT_HTML5, 'UTF-8');

        foreach ($placeholders as $placeholder => $html) {
            $text = str_replace(
                htmlspecialchars($placeholder, ENT_QUOTES | ENT_HTML5, 'UTF-8'),
                $html,
                $text,
            );
        }

        return $text;
    }

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders): string
    {
        for ($pass = 0; $pass < self::MAX_PASSES; $pass++) {
            foreach ($this->converters as $converter) {
                $text = $converter->convert($text, $placeholders, $this);
            }
        }

        return $text;
    }
}
