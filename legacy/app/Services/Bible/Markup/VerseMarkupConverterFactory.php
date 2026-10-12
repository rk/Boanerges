<?php

namespace App\Services\Bible\Markup;

final class VerseMarkupConverterFactory
{
    /**
     * @return list<VerseMarkupConverter>
     */
    public static function defaultConverters(): array
    {
        return [
            new RemoveVerseMarkupConverter([
                // OSIS notes and reference material
                '/<note\b[^>]*>.*?<\/note>/is',
                '/<figure\b[^>]*>.*?<\/figure>/is',
                '/<title\b[^>]*>.*?<\/title>/is',
                '/<rdgGroup\b[^>]*>.*?<\/rdgGroup>/is',
                '/<rdg\b[^>]*>.*?<\/rdg>/is',
                '/<catchWord\b[^>]*>.*?<\/catchWord>/is',
                '/<index\b[^>]*>.*?<\/index>/is',
                '/<milestone\b[^>]*\/>/is',
                '/<milestone\b[^>]*>.*?<\/milestone>/is',
                // GBF footnotes and comments (GBF close tags omit the slash, e.g. <Fn> not </Fn>)
                '/<FN[^>]*>.*?<Fn>/s',
                '/<RF[^>]*>.*?<Rf>/s',
                '/<CM>.*?<Cm>/s',
            ]),
            new OsisHiVerseMarkupConverter(),
            new OsisSemanticVerseMarkupConverter(),
            // GBF font attributes (e.g. YLT)
            new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
            new PairTagVerseMarkupConverter('FB', 'Fb', 'strong'),
            new PairTagVerseMarkupConverter('FR', 'Fr', 'span', 'words-of-jesus'),
            new PairTagVerseMarkupConverter('FO', 'Fo', 'span', 'ot-quote'),
            new PairTagVerseMarkupConverter('FS', 'Fs', 'sup'),
            new PairTagVerseMarkupConverter('FU', 'Fu', 'span', 'underline'),
            new UnwrapVerseMarkupConverter([
                'p', 'l', 'lg', 'q', 'a', 'w', 'inscription', 'mentioned', 'name',
                'reference', 'seg', 'salute', 'signed', 'closer', 'speech', 'speaker',
                'list', 'item', 'table', 'head', 'row', 'cell', 'caption', 'chapter',
                'div', 'verse', 'abbr', 'hi', 'transChange', 'divineName', 'foreign',
            ]),
            new RemoveVerseMarkupConverter([
                // GBF structural and metadata tags
                '/<RX[^>]*>.*?<\/Rx>/s',
                '/<RP[^>]*>/s',
                '/<Rp[^>]*>/s',
                '/<H[^>]*>/s',
                '/<B[^>]*>/s',
                '/<D[^>]*>/s',
                '/<J[^>]*>/s',
                '/<P.>/s',
                '/<W[^>]*>/s',
                '/<S[^>]*>/s',
                '/<N[^>]*>/s',
                '/<C.>/s',
                '/<TS[^>]*>.*?<\/Ts>/s',
            ]),
        ];
    }

    public static function defaultFormatter(): VerseTextFormatter
    {
        return new VerseTextFormatter(self::defaultConverters());
    }
}
