<?php

namespace App\Services\Bible\Markup;

final class RemoveVerseMarkupConverter implements VerseMarkupConverter
{
    /**
     * @param  list<string>  $patterns
     */
    public function __construct(
        private array $patterns,
    ) {}

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders, ?VerseTextFormatter $formatter = null): string
    {
        foreach ($this->patterns as $pattern) {
            $text = (string) preg_replace($pattern, '', $text);
        }

        return $text;
    }
}
