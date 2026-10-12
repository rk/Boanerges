<?php

namespace App\Services\Bible\Markup;

final class UnwrapVerseMarkupConverter implements VerseMarkupConverter
{
    /**
     * @param  list<string>  $tagNames
     */
    public function __construct(
        private array $tagNames,
    ) {}

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders, ?VerseTextFormatter $formatter = null): string
    {
        foreach ($this->tagNames as $tagName) {
            $quotedTag = preg_quote($tagName, '/');

            $text = (string) preg_replace('/<' . $quotedTag . '\b[^>]*>/i', '', $text);
            $text = (string) preg_replace('/<\/' . $quotedTag . '>/i', '', $text);
        }

        return $text;
    }
}
