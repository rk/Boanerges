<?php

namespace Native\Desktop\Facades;

use Native\Desktop\DataObjects\Printer;

/**
 * @method static list<Printer> printers()
 * @method static void print(string $html, ?Printer $printer = null, array<string, mixed> $settings = [])
 * @method static string printToPDF(string $html, array<string, mixed> $settings = [])
 */
class System extends \Illuminate\Support\Facades\Facade {}
