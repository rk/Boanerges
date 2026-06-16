<?php

namespace App\Data;

readonly class CatalogEntry
{

    public function __construct(
        public string $short,
        public string $name,
        public string $url,
        public ?string $about = null,
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
     * @param  array{short: string, name: string, url: string, about: string|null}  $attributes
     */
    public static function fromArray(array $attributes): self
    {
        return new self(
            short: $attributes['short'],
            name: $attributes['name'],
            url: $attributes['url'] ?? 'https://crosswire.org/ftpmirror/pub/sword/packages/rawzip/' . $attributes['short'],
            about: $attributes['about'] ?? 'https://crosswire.org/sword/modules/ModInfo.jsp?modName=' . $attributes['short'],
        );
    }
}
