<?php

namespace App\Services\Bible\Import;

trait ProgressEmitter
{
    protected int $progressBase;

    protected int $progressAmount;

    protected mixed $progressCb;

    /**
     * @param int $base The base progress this pads by.
     * @param int $amount The amount this class's process amounts for.
     * @param callable $cb The emitter callback that receives a float percentage as its first argument.
     */
    public function progressConfigure(int $base, int $amount, mixed $cb): void
    {
        $this->progressBase = $base;
        $this->progressAmount = $amount;
        $this->progressCb = $cb;
    }

    public function onProgress(int $num, int $of): void
    {
        $cb = $this->progressCb;
        $cb($num / $of * $this->progressAmount + $this->progressBase);
    }
}
