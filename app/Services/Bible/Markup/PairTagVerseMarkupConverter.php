<?php

namespace App\Services\Bible\Markup;

final class PairTagVerseMarkupConverter implements VerseMarkupConverter
{
    use CreatesVerseMarkupPlaceholders;

    public function __construct(
        private string $openTag,
        private string $closeTag,
        private string $htmlTag,
        private ?string $class = null,
    ) {}

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders, ?VerseTextFormatter $formatter = null): string
    {
        $pattern = '/<'
            .preg_quote($this->openTag, '/')
            .'>(.*?)<'
            .preg_quote($this->closeTag, '/')
            .'>/s';

        return (string) preg_replace_callback(
            $pattern,
            function (array $matches) use (&$placeholders, $formatter): string {
                $inner = $formatter !== null
                    ? $formatter->convert($matches[1], $placeholders)
                    : $matches[1];

                return $this->createPlaceholder(
                    $this->wrapContent($this->htmlTag, $inner, $this->class, escape: $formatter === null),
                    $placeholders,
                );
            },
            $text,
        );
    }
}
