<?php

namespace App\Exceptions;

use RuntimeException;

class BibleModuleNotInstalledException extends RuntimeException
{
    public static function missing(string $moduleKey): self
    {
        return new self(
            "The SWORD module \"{$moduleKey}\" is not installed. Run `php artisan bible:install-asv`.",
        );
    }
}
