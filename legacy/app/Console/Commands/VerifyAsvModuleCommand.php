<?php

namespace App\Console\Commands;

use App\Services\Bible\BibleModuleManager;
use Illuminate\Console\Command;

class VerifyAsvModuleCommand extends Command
{
    protected $signature = 'bible:verify-asv';

    protected $description = 'Verify the bundled ASV SWORD module is readable';

    public function handle(BibleModuleManager $manager): int
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
}
