<?php

namespace App\Data;

use App\Enums\TranslationInstallStatus;
use App\Enums\TranslationInstallStep;

readonly class TranslationConfig
{
    public function __construct(
        public string $id,
        public string $module,
        public string $name,
        public string $abbrev,
        public bool $bundled = false,
        public ?string $about = null,
        public ?TranslationInstallStatus $installStatus = null,
        public ?TranslationInstallStep $installStep = null,
        public ?string $installError = null,
    ) {}

    /**
     * @param  array{id: string, module: string, name: string, abbrev: string, bundled?: bool, about?: string|null, installStatus?: string|null, installStep?: string|null, installError?: string|null}  $attributes
     */
    public static function fromArray(array $attributes): self
    {
        return new self(
            id: $attributes['id'],
            module: $attributes['module'],
            name: $attributes['name'],
            abbrev: $attributes['abbrev'],
            bundled: $attributes['bundled'] ?? false,
            about: $attributes['about'] ?? null,
            installStatus: isset($attributes['installStatus'])
                ? TranslationInstallStatus::tryFrom((string) $attributes['installStatus'])
                : null,
            installStep: isset($attributes['installStep'])
                ? TranslationInstallStep::tryFrom((string) $attributes['installStep'])
                : null,
            installError: $attributes['installError'] ?? null,
        );
    }
}
