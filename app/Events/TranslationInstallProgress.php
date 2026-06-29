<?php

namespace App\Events;

use App\Enums\TranslationInstallStep;
use Illuminate\Broadcasting\Channel;
use Illuminate\Contracts\Broadcasting\ShouldBroadcastNow;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class TranslationInstallProgress implements ShouldBroadcastNow
{
    use Dispatchable;
    use SerializesModels;

    public function __construct(
        public string $abbrev,
        public TranslationInstallStep $step,
        public int $percent,
        public ?string $error = null,
    ) {}

    /**
     * @return array<string, mixed>
     */
    public function broadcastWith(): array
    {
        return [
            'abbrev' => $this->abbrev,
            'step' => $this->step->value,
            'percent' => $this->percent,
            'error' => $this->error,
        ];
    }

    /**
     * @return array<int, Channel>
     */
    public function broadcastOn(): array
    {
        return [new Channel('nativephp')];
    }
}
