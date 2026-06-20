<?php

namespace App\Services;

class BatchedInsertQueue
{
    protected array $rows = [];

    /** @param callable $process */
    public function __construct(
        protected readonly mixed $process,
        protected readonly int $batchSize = 500,
    ) {}

    public function push(array $row): void
    {
        $this->rows[] = $row;

        if (count($this->rows) >= $this->batchSize) {
            $this->processRows();
        }
    }

    public function done(): void
    {
        if (! empty($this->rows)) {
            $this->processRows();
        }
    }

    protected function processRows(): void
    {
        $fn = $this->process;

        $fn($this->rows);

        $this->rows = [];
    }
}
