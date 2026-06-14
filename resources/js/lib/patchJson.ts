import { http } from '@inertiajs/svelte';

export async function patchJson<T>(url: string, data: Record<string, unknown>): Promise<T> {
    const response = await http.getClient().request({
        method: 'patch',
        url,
        data: JSON.stringify(data),
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
    });

    if (response.status < 200 || response.status >= 300) {
        throw new Error(`Request failed: ${response.status}`);
    }

    return JSON.parse(response.data) as T;
}
