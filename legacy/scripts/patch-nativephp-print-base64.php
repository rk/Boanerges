<?php

$root = dirname(__DIR__);

$files = [
    $root . '/vendor/nativephp/desktop/src/System.php' => [
        [
            "        \$this->client->post('system/print', [\n            'html' => \$html,",
            "        \$this->client->post('system/print', [\n            'html' => base64_encode(\$html),",
        ],
    ],
    $root . '/vendor/nativephp/desktop/resources/electron/electron-plugin/src/server/api/system.ts' => [
        [
            "    await printWindow.loadURL(`data:text/html;charset=UTF-8,\${html}`);\n});\n\nrouter.post('/print-to-pdf'",
            "    await printWindow.loadURL(`data:text/html;base64;charset=UTF-8,\${html}`);\n});\n\nrouter.post('/print-to-pdf'",
        ],
    ],
    $root . '/vendor/nativephp/desktop/resources/electron/electron-plugin/dist/server/api/system.js' => [
        [
            "    yield printWindow.loadURL(`data:text/html;charset=UTF-8,\${html}`);\n}));\nrouter.post('/print-to-pdf'",
            "    yield printWindow.loadURL(`data:text/html;base64;charset=UTF-8,\${html}`);\n}));\nrouter.post('/print-to-pdf'",
        ],
    ],
];

foreach ($files as $path => $replacements) {
    if (! is_file($path)) {
        continue;
    }

    $content = file_get_contents($path);

    foreach ($replacements as [$search, $replace]) {
        if (str_contains($content, $search)) {
            $content = str_replace($search, $replace, $content);
        }
    }

    file_put_contents($path, $content);
}
