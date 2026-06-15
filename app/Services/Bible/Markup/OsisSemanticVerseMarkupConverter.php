<?php

namespace App\Services\Bible\Markup;

final class OsisSemanticVerseMarkupConverter implements VerseMarkupConverter
{
    use CreatesVerseMarkupPlaceholders;

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders, ?VerseTextFormatter $formatter = null): string
    {
        $text = (string) preg_replace_callback(
            '/<transChange\b[^>]*\btype=(["\'])([^"\']+)\1[^>]*>(.*?)<\/transChange>/is',
            function (array $matches) use (&$placeholders, $formatter): string {
                $type = strtolower($matches[2]);

                [$htmlTag, $class] = match ($type) {
                    'added' => ['em', null],
                    'removed' => ['del', null],
                    default => ['span', 'trans-change'],
                };

                $inner = $formatter !== null
                    ? $formatter->convert($matches[3], $placeholders)
                    : $matches[3];

                return $this->createPlaceholder(
                    $this->wrapContent($htmlTag, $inner, $class, escape: $formatter === null),
                    $placeholders,
                );
            },
            $text,
        );

        $text = (string) preg_replace_callback(
            '/<divineName>(.*?)<\/divineName>/is',
            function (array $matches) use (&$placeholders, $formatter): string {
                $inner = $formatter !== null
                    ? $formatter->convert($matches[1], $placeholders)
                    : $matches[1];

                return $this->createPlaceholder(
                    $this->wrapContent('span', $inner, 'divine-name', escape: $formatter === null),
                    $placeholders,
                );
            },
            $text,
        );

        return (string) preg_replace_callback(
            '/<foreign>(.*?)<\/foreign>/is',
            function (array $matches) use (&$placeholders, $formatter): string {
                $inner = $formatter !== null
                    ? $formatter->convert($matches[1], $placeholders)
                    : $matches[1];

                return $this->createPlaceholder(
                    $this->wrapContent('span', $inner, 'foreign', escape: $formatter === null),
                    $placeholders,
                );
            },
            $text,
        );
    }
}
