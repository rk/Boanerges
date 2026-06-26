import { createInertiaApp } from '@inertiajs/svelte';
import { registerApplicationMenuHandlers } from '@/lib/applicationMenu.ts';

const appName = import.meta.env.VITE_APP_NAME || 'Laravel';

registerApplicationMenuHandlers();

createInertiaApp({
    title: (title) => (title ? `${title} - ${appName}` : appName),
    progress: {
        color: '#4B5563',
    },
});
