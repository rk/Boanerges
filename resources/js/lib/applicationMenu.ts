import { openTranslationManager } from '@/lib/bible.svelte.ts';
import { onNativeEvent } from '@/lib/nativeBroadcast.ts';
import { openPrintOptions } from '@/lib/printStudy.svelte.ts';
import {
    ensureCrossReferencesColumn,
    ensureSearchColumn,
    openSettings,
    setColumnCount,
    setScrollSync,
} from '@/lib/study.svelte.ts';

const MENU_EVENT = 'Native\\Desktop\\Events\\Menu\\MenuItemClicked';

export function registerApplicationMenuHandlers(): void {
    onNativeEvent<{ item?: { id?: string; checked?: boolean } }>(MENU_EVENT, (payload) => {
        const id = payload.item?.id;

        if (! id) {
            return;
        }

        switch (id) {
            case 'study.search':
                ensureSearchColumn();
                break;
            case 'study.cross-references':
                ensureCrossReferencesColumn();
                break;
            case 'study.print':
                openPrintOptions();
                break;
            case 'file.manage-translations':
                openTranslationManager();
                break;
            case 'view.settings':
                openSettings();
                break;
            case 'view.scroll-sync':
                setScrollSync(Boolean(payload.item?.checked));
                break;
            case 'view.columns.1':
                setColumnCount(1);
                break;
            case 'view.columns.2':
                setColumnCount(2);
                break;
            case 'view.columns.3':
                setColumnCount(3);
                break;
            default:
                break;
        }
    });
}
