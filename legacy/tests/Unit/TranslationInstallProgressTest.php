<?php

use App\Enums\TranslationInstallStatus;
use App\Enums\TranslationInstallStep;

test('every install status returns a progress percent or null for pending fallback', function () {
    foreach (TranslationInstallStatus::cases() as $status) {
        $percent = $status->progressPercent();

        if ($percent !== null) {
            expect($percent)->toBeGreaterThanOrEqual(0)
                ->toBeLessThanOrEqual(100);
        }
    }
});

test('every install step returns a progress percent', function () {
    foreach (TranslationInstallStep::cases() as $step) {
        expect($step->progressPercent())->toBeGreaterThanOrEqual(0)
            ->toBeLessThanOrEqual(100);
    }
});

test('active install statuses are flagged', function () {
    expect(TranslationInstallStatus::Importing->isActive())->toBeTrue()
        ->and(TranslationInstallStatus::Ready->isActive())->toBeFalse()
        ->and(TranslationInstallStatus::Failed->isActive())->toBeFalse();
});
