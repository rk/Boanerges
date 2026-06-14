<?php

namespace App\Services\Bible;

use App\Exceptions\BibleModuleNotInstalledException;
use Illuminate\Support\Facades\Storage;
use rk\PhpSword\SwordBible;
use rk\PhpSword\SwordModules;

class BibleModuleManager
{
    /** @var array<string, SwordBible> */
    private array $openModules = [];

    private ?SwordModules $swordModules = null;

    public function __construct(
        private InstalledTranslationRegistry $registry,
    ) {}

    public function swordRoot(): string
    {
        return Storage::disk('extras')->path(config('boanerges.sword_path'));
    }

    public function isInstalled(): bool
    {
        return is_dir($this->swordRoot() . DIRECTORY_SEPARATOR . 'mods.d');
    }

    public function open(string $moduleKey): SwordBible
    {
        if (isset($this->openModules[$moduleKey])) {
            return $this->openModules[$moduleKey];
        }

        if (! $this->isInstalled()) {
            throw BibleModuleNotInstalledException::missing($moduleKey);
        }

        $modules = $this->modules();

        if (! array_key_exists($moduleKey, $modules->modules)) {
            throw BibleModuleNotInstalledException::missing($moduleKey);
        }

        $this->openModules[$moduleKey] = $modules->getBibleFromModule($moduleKey);

        return $this->openModules[$moduleKey];
    }

    public function openForTranslation(string $translationId): SwordBible
    {
        $translation = $this->registry->find($translationId);

        return $this->open($translation->module);
    }

    private function modules(): SwordModules
    {
        if ($this->swordModules === null) {
            $this->swordModules = new SwordModules($this->swordRoot());
            $this->swordModules->parseModules();
        }

        return $this->swordModules;
    }
}
