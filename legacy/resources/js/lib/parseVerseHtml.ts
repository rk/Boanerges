export type VerseHtmlTextNode = { type: 'text'; value: string };

export type VerseHtmlElementNode = {
    type: 'tag';
    name: 'em' | 'strong' | 'sup' | 'span';
    className?: string;
    children: VerseHtmlNode[];
};

export type VerseHtmlNode = VerseHtmlTextNode | VerseHtmlElementNode;

const ALLOWED_TAGS = new Set(['em', 'strong', 'sup', 'span']);

const ALLOWED_SPAN_CLASSES = new Set([
    'words-of-jesus',
    'ot-quote',
    'underline',
    'small-caps',
    'divine-name',
    'foreign',
]);

export function parseVerseHtml(html: string): VerseHtmlNode[] {
    if (html === '') {
        return [];
    }

    const doc = new DOMParser().parseFromString(
        `<body>${html}</body>`,
        'text/html',
    );

    return walkNodes(doc.body.childNodes);
}

function walkNodes(nodes: NodeListOf<ChildNode>): VerseHtmlNode[] {
    const result: VerseHtmlNode[] = [];

    for (const node of nodes) {
        if (node.nodeType === Node.TEXT_NODE) {
            const value = node.textContent ?? '';

            if (value !== '') {
                result.push({ type: 'text', value });
            }

            continue;
        }

        if (node.nodeType !== Node.ELEMENT_NODE) {
            continue;
        }

        const element = node as HTMLElement;
        const tagName = element.tagName.toLowerCase();

        if (!ALLOWED_TAGS.has(tagName)) {
            result.push(...walkNodes(element.childNodes));

            continue;
        }

        if (tagName === 'span') {
            const className = element.className.trim();

            if (className === '' || !ALLOWED_SPAN_CLASSES.has(className)) {
                result.push(...walkNodes(element.childNodes));

                continue;
            }

            result.push({
                type: 'tag',
                name: 'span',
                className,
                children: walkNodes(element.childNodes),
            });

            continue;
        }

        result.push({
            type: 'tag',
            name: tagName as 'em' | 'strong' | 'sup',
            children: walkNodes(element.childNodes),
        });
    }

    return result;
}
