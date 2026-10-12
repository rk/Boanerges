<?php

namespace App\Enums;

enum ReaderTheme: string
{
    case Auto = 'auto';
    case Light = 'light';
    case Autumn = 'autumn';
    case Bumblebee = 'bumblebee';
    case Cupcake = 'cupcake';
    case Emerald = 'emerald';
    case Garden = 'garden';
    case Caramellatte = 'caramellatte';
    case Lemonade = 'lemonade';
    case Nord = 'nord';
    case Pastel = 'pastel';
    case Silk = 'silk';
    case Valentine = 'valentine';
    case Winter = 'winter';
    case Dark = 'dark';
    case Coffee = 'coffee';
    case Dracula = 'dracula';
    case Forest = 'forest';
    case Luxury = 'luxury';
    case Night = 'night';
    case Sunset = 'sunset';
    case Synthwave = 'synthwave';

    public function group(): string
    {
        return match ($this) {
            self::Auto, self::Light, self::Dark => 'basics',
            self::Autumn,
            self::Bumblebee,
            self::Cupcake,
            self::Emerald,
            self::Garden,
            self::Caramellatte,
            self::Lemonade,
            self::Nord,
            self::Pastel,
            self::Silk,
            self::Valentine,
            self::Winter => 'light',
            default => 'dark',
        };
    }

    public function label(): string
    {
        return match ($this) {
            self::Auto => 'Auto',
            self::Light => 'Light',
            self::Dark => 'Dark',
            self::Autumn => 'Autumn',
            self::Bumblebee => 'Bumblebee',
            self::Cupcake => 'Cupcake',
            self::Emerald => 'Emerald',
            self::Garden => 'Garden',
            self::Caramellatte => 'Latte',
            self::Lemonade => 'Lemonade',
            self::Nord => 'Nord',
            self::Pastel => 'Pastel',
            self::Silk => 'Silk',
            self::Valentine => 'Valentine',
            self::Winter => 'Winter',
            self::Coffee => 'Coffee',
            self::Dracula => 'Dracula',
            self::Forest => 'Forest',
            self::Luxury => 'Luxury',
            self::Night => 'Night',
            self::Sunset => 'Sunset',
            self::Synthwave => 'Synthwave',
        };
    }
}
