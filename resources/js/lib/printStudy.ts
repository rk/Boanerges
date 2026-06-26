import { store as printStudyRoute } from '@/actions/App/Http/Controllers/StudyPrintController';
import { study } from '@/lib/study.svelte.ts';
import { normalizeColumns } from '@/lib/studyLayout';

export type PrintMode = 'include-user-work' | 'blank-writing';

export async function printStudy(mode: PrintMode): Promise<void> {
    const response = await fetch(printStudyRoute.url(), {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({
            includeUserWork: mode === 'include-user-work',
            columnCount: study.columnCount,
            columns: normalizeColumns(study.columnCount, study.columns),
            bookId: study.bookId,
            chapter: study.chapter,
            translationId: study.translationId,
            translationBId: study.translationBId,
            translationCId: study.translationCId,
        }),
    });

    if (! response.ok) {
        const body = (await response.json().catch(() => null)) as { message?: string } | null;

        throw new Error(body?.message ?? `Print failed: ${response.status}`);
    }
}
