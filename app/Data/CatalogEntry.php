<?php

namespace App\Data;

readonly class CatalogEntry
{
    /**
     * @param  array{short: string, name: string, url: string}  $attributes
     */
    public function __construct(
        public string $short,
        public string $name,
        public string $url,
    ) {}

    public function id(): string
    {
        return strtolower($this->short);
    }

    public function module(): string
    {
        return $this->short;
    }

    /**
     * @param  array{short: string, name: string, url: string}  $attributes
     */
    public static function fromArray(array $attributes): self
    {
        return new self(
            short: $attributes['short'],
            name: $attributes['name'],
            url: $attributes['url'],
        );
    }
}
