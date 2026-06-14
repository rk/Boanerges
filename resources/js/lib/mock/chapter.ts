import type { Book, Chapter, Translation } from '@/lib/types/bible';

export const translations: Translation[] = [
    { id: 'kjv', name: 'King James Version', abbrev: 'KJV' },
    { id: 'asv', name: 'American Standard Version', abbrev: 'ASV' },
    { id: 'web', name: 'World English Bible', abbrev: 'WEB' },
];

export const books: Book[] = [
    { id: 'gen', name: 'Genesis', abbrev: 'Gen', testament: 'ot', chapters: 50 },
    { id: 'exo', name: 'Exodus', abbrev: 'Exo', testament: 'ot', chapters: 40 },
    { id: 'psa', name: 'Psalms', abbrev: 'Psa', testament: 'ot', chapters: 150 },
    { id: 'isa', name: 'Isaiah', abbrev: 'Isa', testament: 'ot', chapters: 66 },
    { id: 'mat', name: 'Matthew', abbrev: 'Mat', testament: 'nt', chapters: 28 },
    { id: 'mrk', name: 'Mark', abbrev: 'Mrk', testament: 'nt', chapters: 16 },
    { id: 'luk', name: 'Luke', abbrev: 'Luk', testament: 'nt', chapters: 24 },
    { id: 'jhn', name: 'John', abbrev: 'Jhn', testament: 'nt', chapters: 21 },
    { id: 'rom', name: 'Romans', abbrev: 'Rom', testament: 'nt', chapters: 16 },
    { id: 'rev', name: 'Revelation', abbrev: 'Rev', testament: 'nt', chapters: 22 },
];

export const genesis15: Chapter = {
    book: 'Genesis',
    bookAbbrev: 'Gen',
    chapter: 15,
    verses: [
        {
            number: 1,
            text: 'After these things the word of the LORD came unto Abram in a vision, saying, Fear not, Abram: I am thy shield, and thy exceeding great reward.',
            paragraphStart: true,
        },
        {
            number: 2,
            text: 'And Abram said, Lord GOD, what wilt thou give me, seeing I go childless, and the steward of my house is this Eliezer of Damascus?',
        },
        {
            number: 3,
            text: 'And Abram said, Behold, to me thou hast given no seed: and, lo, one born in my house is mine heir.',
        },
        {
            number: 4,
            text: 'And, behold, the word of the LORD came unto him, saying, This shall not be thine heir; but he that shall come forth out of thine own bowels shall be thine heir.',
        },
        {
            number: 5,
            text: 'And he brought him forth abroad, and said, Look now toward heaven, and tell the stars, if thou be able to number them: and he said unto him, So shall thy seed be.',
            paragraphStart: true,
        },
        {
            number: 6,
            text: 'And he believed in the LORD; and he counted it to him for righteousness.',
        },
        {
            number: 7,
            text: 'And he said unto him, I am the LORD that brought thee out of Ur of the Chaldees, to give thee this land to inherit it.',
            paragraphStart: true,
        },
        {
            number: 8,
            text: 'And he said, Lord GOD, whereby shall I know that I shall inherit it?',
        },
        {
            number: 9,
            text: 'And he said unto him, Take me an heifer of three years old, and a she goat of three years old, and a ram of three years old, and a turtledove, and a young pigeon.',
        },
        {
            number: 10,
            text: 'And he took unto him all these, and divided them in the midst, and laid each piece one against another: but the birds divided he not.',
        },
        {
            number: 11,
            text: 'And when the fowls came down upon the carcases, Abram drove them away.',
            paragraphStart: true,
        },
        {
            number: 12,
            text: 'And when the sun was going down, a deep sleep fell upon Abram; and, lo, an horror of great darkness fell upon him.',
        },
        {
            number: 13,
            text: 'And he said unto Abram, Know of a surety that thy seed shall be a stranger in a land that is not theirs, and shall serve them; and they shall afflict them four hundred years;',
        },
        {
            number: 14,
            text: 'And also that nation, whom they shall serve, will I judge: and afterward shall they come out with great substance.',
        },
        {
            number: 15,
            text: 'And thou shalt go to thy fathers in peace; thou shalt be buried in a good old age.',
        },
        {
            number: 16,
            text: 'But in the fourth generation they shall come hither again: for the iniquity of the Amorites is not yet full.',
        },
        {
            number: 17,
            text: 'And it came to pass, that, when the sun went down, and it was dark, behold a smoking furnace, and a burning lamp that passed between those pieces.',
            paragraphStart: true,
        },
        {
            number: 18,
            text: 'In the same day the LORD made a covenant with Abram, saying, Unto thy seed have I given this land, from the river of Egypt unto the great river, the river Euphrates:',
        },
        {
            number: 19,
            text: 'The Kenites, and the Kenizzites, and the Kadmonites,',
        },
        {
            number: 20,
            text: 'And the Hittites, and the Perizzites, and the Rephaims,',
        },
        {
            number: 21,
            text: 'And the Amorites, and the Canaanites, and the Girgashites, and the Jebusites.',
        },
    ],
};

export function getChapter(bookId: string, chapter: number): Chapter {
    const book = books.find((item) => item.id === bookId);

    if (book?.id === 'gen' && chapter === 15) {
        return genesis15;
    }

    return {
        book: book?.name ?? 'Genesis',
        bookAbbrev: book?.abbrev ?? 'Gen',
        chapter,
        verses: [
            {
                number: 1,
                text: 'Placeholder chapter content. Wire SWORD integration to load real text.',
                paragraphStart: true,
            },
        ],
    };
}

export function getAdjacentChapter(
    bookId: string,
    chapter: number,
    direction: 'prev' | 'next',
): { bookId: string; chapter: number } | null {
    const book = books.find((item) => item.id === bookId);

    if (! book) {
        return null;
    }

    if (direction === 'prev') {
        if (chapter > 1) {
            return { bookId, chapter: chapter - 1 };
        }

        const index = books.findIndex((item) => item.id === bookId);

        if (index > 0) {
            const previous = books[index - 1];

            return { bookId: previous.id, chapter: previous.chapters };
        }

        return null;
    }

    if (chapter < book.chapters) {
        return { bookId, chapter: chapter + 1 };
    }

    const index = books.findIndex((item) => item.id === bookId);

    if (index < books.length - 1) {
        const next = books[index + 1];

        return { bookId: next.id, chapter: 1 };
    }

    return null;
}
