<?php

use App\Services\Bible\Import\SwordConfReader;
use Tests\TestCase;

uses(TestCase::class);

it('parses metadata from asv sword conf', function () {
    $path = base_path('extras/sword/mods.d/asv.conf');

    if (! is_file($path)) {
        $this->markTestSkipped('Bundled ASV conf not present.');
    }

    $reader = app(SwordConfReader::class);
    $metadata = $reader->parseFile($path);

    expect($metadata['versification'])->toBe('KJV')
        ->and($metadata['format'])->toBe('osis')
        ->and($metadata['version_string'])->toBe('2.0')
        ->and($metadata['version_date'])->toBe('2021-02-18')
        ->and($metadata['copyright'])->toBe('Public Domain')
        ->and($metadata['source'])->toBe('http://www.ebible.org/bible/asv/')
        ->and($metadata['about'])->toContain('American Standard Version')
        ->and($metadata['about'])->not->toContain('\\par');
});
