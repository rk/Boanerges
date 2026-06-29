<?php

namespace App\Data;

readonly class CatalogEntry
{
    /** @param  'sword'|'usfm'|'accordance'  $importAs */
    public function __construct(
        public string $short,
        public string $name,
        public string $url,
        public ?string $about = null,
        public string $importAs = 'sword',
        public ?string $markupFormat = null,
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
        /** @var 'sword'|'usfm'|'accordance' $importAs */
        $importAs = self::resolveImportAs($attributes);

        return new self(
            short: $attributes['short'],
            name: $attributes['name'],
            url: $attributes['url'] ?? "https://crosswire.org/ftpmirror/pub/sword/packages/rawzip/{$attributes['short']}",
            about: $attributes['about'] ?? "https://crosswire.org/sword/modules/ModInfo.jsp?modName={$attributes['short']}",
            importAs: $importAs,
            markupFormat: self::resolveMarkupFormat($attributes),
        );
    }

    /** @param  array<string, mixed>  $attributes
     * @return 'sword'|'usfm'|'accordance'
     */
    private static function resolveImportAs(array $attributes): string
    {
        if (isset($attributes['import_as'])) {
            $importAs = (string) $attributes['import_as'];

            if (in_array($importAs, ['sword', 'usfm', 'accordance'], true)) {
                return $importAs;
            }
        }

        $format = $attributes['format'] ?? null;

        return match ($format) {
            'sword', 'usfm', 'accordance' => $format,
            default => 'sword',
        };
    }

    /** @param  array<string, mixed>  $attributes */
    private static function resolveMarkupFormat(array $attributes): ?string
    {
        if (isset($attributes['markup_format'])) {
            return strtolower((string) $attributes['markup_format']);
        }

        $format = $attributes['format'] ?? null;

        if ($format !== null && ! in_array($format, ['sword', 'usfm', 'accordance'], true)) {
            return strtolower((string) $format);
        }

        return match (self::resolveImportAs($attributes)) {
            'usfm' => 'usfm',
            default => null,
        };
    }
}
