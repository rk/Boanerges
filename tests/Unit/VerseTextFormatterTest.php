<?php

use App\Services\Bible\Markup\PairTagVerseMarkupConverter;
use App\Services\Bible\Markup\VerseMarkupConverterFactory;
use App\Services\Bible\Markup\VerseTextFormatter;

function defaultFormatter(): VerseTextFormatter
{
    return VerseMarkupConverterFactory::defaultFormatter();
}

test('converts gbf italic tags to em elements', function (): void {
    $input = 'Fear not, Abram, I <FI>am<Fi> a shield to thee, thy reward <FI>is<Fi> exceeding great.';

    expect(defaultFormatter()->format($input))->toBe(
        'Fear not, Abram, I <em>am</em> a shield to thee, thy reward <em>is</em> exceeding great.',
    );
});

test('converts gbf bold and red letter tags', function (): void {
    expect(defaultFormatter()->format('<FB>bold<Fb> and <FR>red<Fr>'))
        ->toBe('<strong>bold</strong> and <span class="words-of-jesus">red</span>');
});

test('converts gbf superscript and underline tags', function (): void {
    expect(defaultFormatter()->format('word<FS>1<Fs> and <FU>underlined<Fu>'))
        ->toBe('word<sup>1</sup> and <span class="underline">underlined</span>');
});

test('converts osis hi tags to semantic html', function (): void {
    expect(defaultFormatter()->format('<hi type="italic">added</hi> and <hi type="bold">strong</hi>'))
        ->toBe('<em>added</em> and <strong>strong</strong>');
});

test('converts osis hi small caps and superscript', function (): void {
    expect(defaultFormatter()->format('<hi type="small-caps">Lord</hi> and <hi type="superscript">2</hi>'))
        ->toBe('<span class="small-caps">Lord</span> and <sup>2</sup>');
});

test('converts osis transChange and divineName tags', function (): void {
    $input = 'darkness <transChange type="added">was</transChange> upon the face of the deep. '
        .'<divineName>God</divineName> moved';

    expect(defaultFormatter()->format($input))->toBe(
        'darkness <em>was</em> upon the face of the deep. <span class="divine-name">God</span> moved',
    );
});

test('converts osis foreign text', function (): void {
    expect(defaultFormatter()->format('called <foreign>Bethlehem</foreign>'))
        ->toBe('called <span class="foreign">Bethlehem</span>');
});

test('removes osis notes and unwraps strongs word tags', function (): void {
    $input = 'In the beginning <w lemma="strong:H430">God</w><note>footnote</note> created';

    expect(defaultFormatter()->format($input))->toBe('In the beginning God created');
});

test('removes gbf footnotes and verse milestones', function (): void {
    $input = '<verse sID="Gen.1.1"/>The earth<RF Gen 1:2>ref<Rf> was<FN>note<Fn> void<verse eID="Gen.1.1"/>';

    expect(defaultFormatter()->format($input))->toBe('The earth was void');
});

test('escapes html outside converted markup', function (): void {
    expect(defaultFormatter()->format('Use <script> & <FI>am<Fi> here'))
        ->toBe('Use &lt;script&gt; &amp; <em>am</em> here');
});

test('leaves plain text unchanged', function (): void {
    expect(defaultFormatter()->format('In the beginning God created the heaven and the earth.'))
        ->toBe('In the beginning God created the heaven and the earth.');
});

test('applies custom converters in registration order', function (): void {
    $formatter = new VerseTextFormatter([
        new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
        new PairTagVerseMarkupConverter('FB', 'Fb', 'strong'),
    ]);

    expect($formatter->format('<FB>bold<Fb> and <FI>italic<Fi>'))
        ->toBe('<strong>bold</strong> and <em>italic</em>');
});

test('handles nested gbf markup across passes', function (): void {
    expect(defaultFormatter()->format('<FI><FB>bold italic<Fb><Fi>'))
        ->toBe('<em><strong>bold italic</strong></em>');
});
