<?php

namespace App\Support;

use App\Services\StudySettingsStore;
use Native\Desktop\Facades\Menu;

class ApplicationMenuBuilder
{
    public function __construct(private StudySettingsStore $studySettings) {}

    public function register(): void
    {
        if (! config('nativephp-internal.running', false)) {
            return;
        }

        $study = $this->studySettings->get();
        $columnCount = (int) $study['columnCount'];

        Menu::create(
            Menu::app(),
            Menu::make(
                Menu::label('Manage Translations')->id('file.manage-translations'),
                Menu::separator(),
                Menu::quit(),
            )->label('File'),
            Menu::edit(),
            Menu::make(
                Menu::label('Search')->id('study.search')->hotkey('CmdOrCtrl+F'),
                Menu::label('Cross-References')->id('study.cross-references')->hotkey('CmdOrCtrl+Shift+R'),
                Menu::separator(),
                Menu::label('Print')->id('study.print')->disabled(),
                Menu::label('Feedback')->id('study.feedback')->disabled(),
            )->label('Study'),
            Menu::make(
                Menu::radio('1 Column', $columnCount === 1)->id('view.columns.1'),
                Menu::radio('2 Columns', $columnCount === 2)->id('view.columns.2'),
                Menu::radio('3 Columns', $columnCount === 3)->id('view.columns.3'),
                Menu::separator(),
                Menu::checkbox('Link scroll', false)->id('view.scroll-sync'),
                Menu::separator(),
                Menu::label('Edit Settings')->id('view.settings')->hotkey('CmdOrCtrl+,'),
                Menu::separator(),
                Menu::fullscreen(),
                Menu::devTools(),
            )->label('View'),
            Menu::window(),
        );
    }
}
