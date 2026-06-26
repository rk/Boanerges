type NativeBridge = {
    on: (event: string, callback: (payload: unknown, event?: unknown) => void) => void;
    contextMenu: (items: Array<{ label: string; click?: () => void }>) => void;
};

declare global {
    interface Window {
        Native?: NativeBridge;
    }
}

const listeners = new Map<string, Set<(payload: unknown) => void>>();
let initialized = false;

function registerListener(eventClass: string, callback: (payload: unknown) => void): () => void {
    if (! listeners.has(eventClass)) {
        listeners.set(eventClass, new Set());
    }

    listeners.get(eventClass)!.add(callback);

    return () => listeners.get(eventClass)?.delete(callback);
}

function initNativeBridge(): void {
    if (initialized) {
        return;
    }

    initialized = true;

    window.addEventListener('native:init', () => {
        if (! window.Native) {
            return;
        }

        for (const eventClass of listeners.keys()) {
            window.Native.on(eventClass, (payload) => {
                listeners.get(eventClass)?.forEach((cb) => cb(payload));
            });
        }
    });
}

export function onNativeEvent<T>(
    eventClass: string,
    callback: (payload: T) => void,
): () => void {
    initNativeBridge();

    const wrapped = (payload: unknown) => callback(payload as T);

    return registerListener(eventClass, wrapped);
}

export function pollInstallStatus(
    url: string,
    onUpdate: (status: { step: string; percent: number; install_status: string; install_error?: string | null }) => void,
    intervalMs = 1500,
): () => void {
    let active = true;

    const poll = async () => {
        while (active) {
            try {
                const response = await fetch(url, {
                    headers: { Accept: 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
                });

                if (response.ok) {
                    const data = await response.json();
                    onUpdate(data);
                }
            } catch {
                // ponytail: ignore transient poll failures
            }

            await new Promise((resolve) => setTimeout(resolve, intervalMs));
        }
    };

    void poll();

    return () => {
        active = false;
    };
}

export function watchInstallProgress(
    eventClass: string,
    pollUrl: string,
    onUpdate: (payload: { step: string; percent: number; install_status?: string; install_error?: string | null; abbrev?: string }) => void,
): () => void {
    const stopPoll = pollInstallStatus(pollUrl, onUpdate);
    const stopNative = onNativeEvent(eventClass, onUpdate);

    return () => {
        stopPoll();
        stopNative();
    };
}
