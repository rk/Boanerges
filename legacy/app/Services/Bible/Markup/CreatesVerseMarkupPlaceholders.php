<?php

namespace App\Services\Bible\Markup;

trait CreatesVerseMarkupPlaceholders
{
    /**
     * @param  array<string, string>  $placeholders
     */
    protected function createPlaceholder(string $html, array &$placeholders): string
    {
        $placeholder = "\x7Fboanerges-markup-" . count($placeholders) . "\x7F";
        $placeholders[$placeholder] = $html;

        return $placeholder;
    }

    protected function wrapContent(
        string $htmlTag,
        string $content,
        ?string $class = null,
        bool $escape = true,
    ): string {
        $inner = $escape
            ? htmlspecialchars($content, ENT_QUOTES | ENT_HTML5, 'UTF-8')
            : $content;

        if ($class === null) {
            return '<' . $htmlTag . '>' . $inner . '</' . $htmlTag . '>';
        }

        return '<' . $htmlTag
            . ' class="' . htmlspecialchars($class, ENT_QUOTES | ENT_HTML5, 'UTF-8') . '">'
            . $inner
            . '</' . $htmlTag . '>';
    }
}
