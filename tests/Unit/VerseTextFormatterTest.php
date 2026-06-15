<?php

use App\Services\Bible\Markup\PairTagVerseMarkupConverter;
use App\Services\Bible\Markup\VerseTextFormatter;

test('converts gbf italic tags to em elements', function (): void {
    $formatter = new VerseTextFormatter([
        new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
    ]);

    $input = 'Fear not, Abram, I <FI>am<Fi> a shield to thee, thy reward <FI>is<Fi> exceeding great.';

    expect($formatter->format($input))->toBe(
        'Fear not, Abram, I <em>am</em> a shield to thee, thy reward <em>is</em> exceeding great.',
    );
});

test('escapes html outside converted markup', function (): void {
    $formatter = new VerseTextFormatter([
        new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
    ]);

    expect($formatter->format('Use <script> & <FI>am<Fi> here'))
        ->toBe('Use &lt;script&gt; &amp; <em>am</em> here');
});

test('leaves plain text unchanged', function (): void {
    $formatter = new VerseTextFormatter([
        new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
    ]);

    expect($formatter->format('In the beginning God created the heaven and the earth.'))
        ->toBe('In the beginning God created the heaven and the earth.');
});

test('applies converters in registration order', function (): void {
    $formatter = new VerseTextFormatter([
        new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
        new PairTagVerseMarkupConverter('FB', 'Fb', 'strong'),
    ]);

    expect($formatter->format('<FB>bold<Fb> and <FI>italic<Fi>'))
        ->toBe('<strong>bold</strong> and <em>italic</em>');
});
