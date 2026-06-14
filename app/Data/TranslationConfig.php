<?php

namespace App\Data;

readonly class TranslationConfig
{
    public function __construct(
        public string $id,
        public string $module,
        public string $name,
        public string $abbrev,
        public bool $bundled = false,
    ) {}

    /**
     * @param  array{id: string, module: string, name: string, abbrev: string, bundled?: bool}  $attributes
     */
    public static function fromArray(array $attributes): self
    {
        return new self(
            id: $attributes['id'],
            module: $attributes['module'],
            name: $attributes['name'],
            abbrev: $attributes['abbrev'],
            bundled: $attributes['bundled'] ?? false,
        );
    }
}
