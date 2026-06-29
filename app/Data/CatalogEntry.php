<?php

namespace App\Data;

use App\Enums\CatalogImportFormat;
use App\Enums\VerseMarkupFormat;

readonly class CatalogEntry
{
    public function __construct(
        public string $short,
        public string $name,
        public string $url,
        public ?string $about = null,
        public CatalogImportFormat $importAs = CatalogImportFormat::Sword,
        public ?VerseMarkupFormat $markupFormat = null,
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
     * @param  array{short: string, name: string, url?: string, about?: string|null, import_as?: string, format?: string, markup_format?: string}  $attributes
     */
    public static function fromArray(array $attributes): self
    {
        $importAs = self::resolveImportAs($attributes);

        return new self(
            short: $attributes['short'],
            name: $attributes['name'],
            url: $attributes['url'] ?? "https://crosswire.org/ftpmirror/pub/sword/packages/rawzip/{$attributes['short']}",
            about: $attributes['about'] ?? "https://crosswire.org/sword/modules/ModInfo.jsp?modName={$attributes['short']}",
            importAs: $importAs,
            markupFormat: self::resolveMarkupFormat($attributes, $importAs),
        );
    }

    /** @param  array<string, mixed>  $attributes */
    private static function resolveImportAs(array $attributes): CatalogImportFormat
    {
        if (isset($attributes['import_as'])) {
            $importAs = CatalogImportFormat::tryFrom((string) $attributes['import_as']);

            if ($importAs !== null) {
                return $importAs;
            }
        }

        $format = $attributes['format'] ?? null;

        return match ($format) {
            'sword', 'usfm', 'accordance' => CatalogImportFormat::from($format),
            default => CatalogImportFormat::Sword,
        };
    }

    /** @param  array<string, mixed>  $attributes */
    private static function resolveMarkupFormat(array $attributes, CatalogImportFormat $importAs): ?VerseMarkupFormat
    {
        if (isset($attributes['markup_format'])) {
            return VerseMarkupFormat::tryFrom(strtolower((string) $attributes['markup_format']));
        }

        $format = $attributes['format'] ?? null;

        if ($format !== null && CatalogImportFormat::tryFrom((string) $format) === null) {
            return VerseMarkupFormat::tryFrom(strtolower((string) $format));
        }

        return match ($importAs) {
            CatalogImportFormat::Usfm => VerseMarkupFormat::Usfm,
            default => null,
        };
    }
}
