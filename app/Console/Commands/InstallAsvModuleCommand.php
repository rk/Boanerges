<?php

namespace App\Console\Commands;

use App\Services\Bible\BibleModuleManager;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Storage;
use ZipArchive;

class InstallAsvModuleCommand extends Command
{
    private const string DOWNLOAD_URL = 'https://crosswire.org/sword/servlet/SwordMod.Verify?modName=ASV&pkgType=raw';

    private const string MIRROR_URL = 'http://www.crosswire.org/ftpmirror/pub/sword/packages/rawzip/ASV.zip';

    protected $signature = 'bible:install-asv {--force : Re-download and replace the existing module}';

    protected $description = 'Download and install the ASV SWORD Bible module';

    public function handle(BibleModuleManager $manager): int
    {
        $swordPath = Storage::disk('extras')->path(config('boanerges.sword_path'));
        $confPath = $swordPath . DIRECTORY_SEPARATOR . 'mods.d' . DIRECTORY_SEPARATOR . 'ASV.conf';

        if (is_file($confPath) && ! $this->option('force')) {
            $this->info('ASV module already installed.');

            return $this->verify($manager);
        }

        if (! is_dir($swordPath)) {
            mkdir($swordPath, 0755, true);
        }

        $zipPath = storage_path('app/asv.sword-module.zip');

        $this->info('Downloading ASV module...');

        if (! $this->downloadTo($zipPath)) {
            $this->error('Failed to download ASV module from CrossWire.');

            return self::FAILURE;
        }

        $this->info('Extracting module...');

        $zip = new ZipArchive();

        if ($zip->open($zipPath) !== true) {
            $this->error('Failed to open downloaded archive.');

            return self::FAILURE;
        }

        $zip->extractTo($swordPath);
        $zip->close();

        @unlink($zipPath);

        if (! is_file($confPath)) {
            $this->error('Installation failed: mods.d/ASV.conf not found after extraction.');

            return self::FAILURE;
        }

        $this->info('ASV module installed successfully.');

        return $this->verify($manager);
    }

    private function verify(BibleModuleManager $manager): int
    {
        try {
            $bible = $manager->open('ASV');
            $text = trim($bible->get(
                books: 'Genesis',
                chapters: 1,
                verses: 1,
                clean: true,
                join: '',
            ));

            $this->info("Smoke test passed: {$text}");

            return self::SUCCESS;
        } catch (\Throwable $exception) {
            $this->error($exception->getMessage());

            return self::FAILURE;
        }
    }

    private function downloadTo(string $zipPath): bool
    {
        foreach ([self::DOWNLOAD_URL, self::MIRROR_URL] as $url) {
            $response = Http::timeout(120)->sink($zipPath)->get($url);

            if ($response->successful() && is_file($zipPath) && filesize($zipPath) > 0) {
                return true;
            }
        }

        return false;
    }
}
