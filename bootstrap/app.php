<?php

use App\Http\Middleware\EnsureBundledDataMiddleware;
use App\Http\Middleware\HandleInertiaRequests;
use App\Exceptions\BibleModuleNotInstalledException;
use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Http\Middleware\AddLinkHeadersForPreloadedAssets;
use Illuminate\Http\Request;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__ . '/../routes/web.php',
        commands: __DIR__ . '/../routes/console.php',
        health: '/up',
        then: function (): void {
            require base_path('routes/bible.php');
            require base_path('routes/scribe.php');
        },
    )
    ->withMiddleware(function (Middleware $middleware): void {
        $middleware->web(append: [
            EnsureBundledDataMiddleware::class,
            HandleInertiaRequests::class,
            AddLinkHeadersForPreloadedAssets::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        $exceptions->shouldRenderJsonWhen(
            fn(Request $request) => $request->is('api/*') || $request->expectsJson(),
        );

        $exceptions->render(function (BibleModuleNotInstalledException $exception, Request $request) {
            if ($request->is('bible/*') || $request->expectsJson()) {
                return response()->json(['message' => $exception->getMessage()], 503);
            }
        });
    })->create();
