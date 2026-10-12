<?php

namespace App\Support;

use App\Contracts\BundledContentProvider;
use App\Enums\TranslationInstallStatus;
use App\Models\Translation;
use InvalidArgumentException;

class BundledContentRegistry
{
    /** @var array<string, BundledContentProvider> */
    private array $providers = [];

    public function register(BundledContentProvider $provider): void
    {
        $this->providers[$provider->key()] = $provider;
    }

    public function get(string $key): BundledContentProvider
    {
        $provider = $this->providers[$key] ?? null;

        if ($provider === null) {
            throw new InvalidArgumentException("Unknown bundled content provider: {$key}");
        }

        return $provider;
    }

    /**
     * @return list<BundledContentProvider>
     */
    public function all(): array
    {
        return array_values($this->providers);
    }

    /**
     * @return list<BundledContentProvider>
     */
    public function pending(): array
    {
        $readyTranslationExists = Translation::query()
            ->where('install_status', TranslationInstallStatus::Ready)
            ->exists();

        return array_values(array_filter(
            $this->providers,
            function (BundledContentProvider $provider) use ($readyTranslationExists): bool {
                if ($provider->isImported()) {
                    return false;
                }

                if ($provider->requiresReadyTranslation() && ! $readyTranslationExists) {
                    return false;
                }

                return true;
            },
        ));
    }
}
