<?php

test('returns a successful response', function () {
    $response = $this->get(route('home'));

    $response->assertOk();
});

test('renders the study page', function () {
    $response = $this->get(route('home'));

    $response->assertOk();
    $response->assertInertia(fn($page) => $page->component('Study'));
});

test('scaffolded ui files exist', function () {
    $paths = [
        'resources/js/layouts/AppLayout.svelte',
        'resources/js/pages/Study.svelte',
        'resources/js/views/BibleView.svelte',
        'resources/js/views/ComparisonView.svelte',
        'resources/js/views/ScribeView.svelte',
        'resources/js/components/sidebar/AppSidebar.svelte',
        'resources/js/components/reader/ReaderPane.svelte',
        'resources/js/components/scribe/ScribeEditor.svelte',
        'resources/js/components/settings/ReadabilitySettings.svelte',
    ];

    foreach ($paths as $path) {
        expect(base_path($path))->toBeFile();
    }
});
