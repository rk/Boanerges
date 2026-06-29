import { describe, expect, it } from 'vitest';

import { parseVerseHtml } from '@/lib/parseVerseHtml';
import type { VerseHtmlNode } from '@/lib/parseVerseHtml';

function text(value: string): VerseHtmlNode {
    return { type: 'text', value };
}

function tag(
    name: 'em' | 'strong' | 'sup' | 'span',
    children: VerseHtmlNode[],
    className?: string,
): VerseHtmlNode {
    if (name === 'span') {
        return { type: 'tag', name, className, children };
    }

    return { type: 'tag', name, children };
}

describe('parseVerseHtml', () => {
    it('returns empty array for empty input', () => {
        expect(parseVerseHtml('')).toEqual([]);
    });

    it('parses plain text', () => {
        expect(parseVerseHtml('In the beginning God created.')).toEqual([
            text('In the beginning God created.'),
        ]);
    });

    it('parses em and strong tags', () => {
        expect(
            parseVerseHtml(
                'Fear not, I <em>am</em> and <strong>bold</strong>.',
            ),
        ).toEqual([
            text('Fear not, I '),
            tag('em', [text('am')]),
            text(' and '),
            tag('strong', [text('bold')]),
            text('.'),
        ]);
    });

    it('parses styled span and sup tags', () => {
        expect(
            parseVerseHtml(
                '<span class="words-of-jesus">red</span> and word<sup>1</sup>',
            ),
        ).toEqual([
            tag('span', [text('red')], 'words-of-jesus'),
            text(' and word'),
            tag('sup', [text('1')]),
        ]);
    });

    it('parses nested markup', () => {
        expect(parseVerseHtml('<em><strong>bold italic</strong></em>')).toEqual(
            [tag('em', [tag('strong', [text('bold italic')])])],
        );
    });

    it('flattens disallowed tags to text content', () => {
        expect(
            parseVerseHtml('Use &lt;script&gt; &amp; <em>am</em> here'),
        ).toEqual([
            text('Use <script> & '),
            tag('em', [text('am')]),
            text(' here'),
        ]);
    });

    it('unwraps disallowed span classes', () => {
        expect(parseVerseHtml('<span class="evil">bad</span>')).toEqual([
            text('bad'),
        ]);
    });
});
