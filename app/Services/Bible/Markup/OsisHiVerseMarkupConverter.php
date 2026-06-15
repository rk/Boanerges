<?php

namespace App\Services\Bible\Markup;

final class OsisHiVerseMarkupConverter implements VerseMarkupConverter
{
    use CreatesVerseMarkupPlaceholders;

    /**
     * @var array<string, array{0: string, 1: ?string}>
     */
    private const TYPE_MAP = [
        'italic' => ['em', null],
        'bold' => ['strong', null],
        'small-caps' => ['span', 'small-caps'],
        'superscript' => ['sup', null],
        'subscript' => ['sub', null],
        'underline' => ['span', 'underline'],
        'line-through' => ['del', null],
    ];

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders, ?VerseTextFormatter $formatter = null): string
    {
        return (string) preg_replace_callback(
            '/<hi\b[^>]*\btype=(["\'])([^"\']+)\1[^>]*>(.*?)<\/hi>/is',
            function (array $matches) use (&$placeholders, $formatter): string {
                $type = strtolower($matches[2]);
                $mapping = self::TYPE_MAP[$type] ?? null;

                if ($mapping === null) {
                    return $matches[3];
                }

                [$htmlTag, $class] = $mapping;
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
    }
}
