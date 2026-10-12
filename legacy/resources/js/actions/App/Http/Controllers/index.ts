import SettingsController from './SettingsController'
import BibleController from './BibleController'
import DictionaryController from './DictionaryController'
import ScribeController from './ScribeController'
import NotesController from './NotesController'
import StudyPrintController from './StudyPrintController'
import VerseListController from './VerseListController'

const Controllers = {
    SettingsController: Object.assign(SettingsController, SettingsController),
    BibleController: Object.assign(BibleController, BibleController),
    DictionaryController: Object.assign(DictionaryController, DictionaryController),
    ScribeController: Object.assign(ScribeController, ScribeController),
    NotesController: Object.assign(NotesController, NotesController),
    StudyPrintController: Object.assign(StudyPrintController, StudyPrintController),
    VerseListController: Object.assign(VerseListController, VerseListController),
}

export default Controllers