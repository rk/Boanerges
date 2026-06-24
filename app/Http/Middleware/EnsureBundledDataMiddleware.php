<?php

namespace App\Http\Middleware;

use App\Listeners\EnsureBundledData;
use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Schema;
use Symfony\Component\HttpFoundation\Response;
use Throwable;

class EnsureBundledDataMiddleware
{
    public function __construct(
        private EnsureBundledData $ensureBundledData,
    ) {}

    public function handle(Request $request, Closure $next): Response
    {
        if (Schema::hasTable('translations')) {
            try {
                $this->ensureBundledData->handle();
            } catch (Throwable) {
                // ponytail: allow requests through while bootstrap imports recover
            }
        }

        return $next($request);
    }
}
