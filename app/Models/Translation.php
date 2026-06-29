<?php

namespace App\Models;

use App\Enums\TranslationInstallStatus;
use Illuminate\Database\Eloquent\Model;

/**
 * @property string $abbrev
 * @property string $name
 * @property string|null $format
 * @property string|null $install_step
 * @property string|null $install_error
 * @property bool $bundled
 * @property TranslationInstallStatus $install_status
 */
class Translation extends Model
{
    /** @var list<string> */
    protected $fillable = [
        'abbrev',
        'name',
        'format',
        'versification',
        'about',
        'version_string',
        'version_date',
        'copyright',
        'copyright_contact',
        'source',
        'install_status',
        'install_step',
        'install_error',
        'bundled',
    ];

    /** @return array<string, string> */
    protected function casts(): array
    {
        return [
            'install_status' => TranslationInstallStatus::class,
            'bundled' => 'boolean',
        ];
    }

    public function isReady(): bool
    {
        return $this->install_status === TranslationInstallStatus::Ready;
    }

    public function installStatusValue(): string
    {
        $status = $this->install_status;

        return $status instanceof TranslationInstallStatus ? $status->value : (string) $status;
    }

    public function updateProgress(TranslationInstallStatus $status, string $step, int $percent): void
    {
        if ($this->install_status !== $status || $this->install_step !== $step) {
            $this->update([
                'install_status' => $status,
                'install_step' => $step,
                'install_error' => null,
            ]);
        }

        event(new \App\Events\TranslationInstallProgress(
            abbrev: $this->abbrev,
            step: $step,
            percent: $percent,
        ));
    }

    public function markFailed(string $error): void
    {
        $this->update([
            'install_status' => TranslationInstallStatus::Failed,
            'install_error' => $error,
        ]);

        event(new \App\Events\TranslationInstallProgress(
            abbrev: $this->abbrev,
            step: 'failed',
            percent: 0,
            error: $error,
        ));
    }
}
