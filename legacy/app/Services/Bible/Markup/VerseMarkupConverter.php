<?php

namespace App\Services\Bible\Markup;

interface VerseMarkupConverter
{
    /**
     * @param  array<string, string>  $placeholders
     */
    public function convert(string $text, array &$placeholders, ?VerseTextFormatter $formatter = null): string;
}
