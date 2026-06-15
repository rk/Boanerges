<?php

namespace App\Services\Bible\Markup;

final class PairTagVerseMarkupConverter implements VerseMarkupConverter
{
    public function __construct(
        private string $openTag,
        private string $closeTag,
        private string $htmlTag,
    ) {}

    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders): string
    {
        $pattern = '/<'
            .preg_quote($this->openTag, '/')
            .'>(.*?)<'
            .preg_quote($this->closeTag, '/')
            .'>/s';

        return (string) preg_replace_callback(
            $pattern,
            function (array $matches) use (&$placeholders): string {
                $placeholder = "\x7Fboanerges-markup-".count($placeholders)."\x7F";
                $placeholders[$placeholder] = '<'.$this->htmlTag.'>'
                    .htmlspecialchars($matches[1], ENT_QUOTES | ENT_HTML5, 'UTF-8')
                    .'</'.$this->htmlTag.'>';

                return $placeholder;
            },
            $text,
        );
    }
}
