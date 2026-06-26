# Boanerges Design

- Our front-end design will be via DaisyUI, using the DaisyUI skill.
- We will introduce each app component as a Svelte component.
- Lucide for icons.

## Revised Layout

- A "splash" page is shown on first launch, during which time the Bibles bundled with the app are imported and installed. When 1+ translation is installed, this screen no longer displays.
- The main layout of the app is the Bible view. It comprises several sections:
  - The toolbar. It shows the brand, followed by actions that may be taken related to the user's current view.
  - The view (main region), which should be highly optimized for readability.
- There will be an Edit Settings dialog to allow control of font size, line height, and font family. It will also allow control over the theme.

## The Toolbar

Toolbar zones (left to right): brand, book/chapter picker, manage translations, column count (1/2/3), scroll-sync toggle when multiple Bible columns are visible, and Edit Settings.

Each column has its own header bar with a view selector (columns 2–3 only; column 1 is locked to Translation) and icon-led controls: Book icon + translation select for Bible columns, chapter heading + save icon for Notes/Scribe, search input for Search, and verse reference input for Cross References.

Study actions (Search, Cross-References, Print, Feedback) live in the native **Study** menu. File/View menus cover translation management, column layout, scroll sync, and settings. Cross-References is a column type; right-click a verse in any Bible column to open cross references in the last secondary slot (or the slot already showing cross references).

The goal of the toolbar is improved UX by reducing user friction when:

- Setting to single-column focus mode. It should display the primary translation with no distractions or other columns.
- Setting to two-column mode, with free selection between what shows in the right-hand column: the secondary translation, notes, or the scribe pane.
- Setting to three-column mode, with free selection between what shows in the two additional columns, as in two-column mode. (However, no two columns may show the same thing.)

Persisted layout state replaces the legacy `activeView`:

- `columnCount`: 1, 2, or 3
- `columns`: additional column types only (`bible-secondary`, `notes`, `scribe`, `search`, `cross-references`); column 0 is always the primary Bible (`translationId`)
- `translationBId` / `translationCId`: translations for secondary Bible slots

## The column displays / views

When the primary translation navigates, all other columns should also navigate to the same chapter. The chapter headings should stay the same as before. The vertical navigation between chapters by pre- and post- elements above/below the text (but inline scrolling) is acceptable.

### Bible View

This is just a display of the Bible in the column, either primary translation or secondary translation. Bible columns use paragraph layout with inline verse numbers.

### Notes View

This view displays a free text editor to the user that is saved for the chapter, unified across all translations.

### Scribe View

This view displays a verse-by-verse editable chapter, numbered and with paragraph breaks based upon the current primary translation.

### Search View

A way to search in a column, and open results in all open Bible views. Search box at the top, paginated results inside. When a search result is tapped and navigates the Bible views to that reference, the same highlight as the cross-reference uses shall be applied.

### Cross References View

A column showing cross references for a verse reference entered in the header. Results navigate all Bible columns and apply verse highlight. Open via Study menu, or right-click any verse in a Bible column.

# Legacy Design

## Layouts

- There will always be a left sidebar used to: change views, change translation, change books (in the translation), and change chapter within the current book.
  - The view selector shall be a toggle between single (default / "Bible"), Comparison, and Scribe.
  - Book selection will be a left/right column for Old and New Testaments (ot/nt).
- The default layout will be the Bible view, where the main pane is a single column paragraph Bible reader.
- Two alternate layouts will be: Comparison view, a two column view of the same book/chapter in two different translations; and Scribe view, a three column view with two different translations to the left/right of an editor.
- There should be menu item that allows the user to edit settings, probably to allow control over font-size, line-height, as well as color/contrast.
- The sidebar should probably be collapsible.

### Bible View

The main pane of the Bible view is a single column, centered, vertically scrolling chapter of the Bible. Text will be formatted as paragraphs with verse numbers inline.

Before the top of the chapter will be a divider like "^ Genesis 14 ^" that the user can click to go to the previous chapter. A similar divider like "v Genesis 16 v" will be clickable to advance to the next chapter. At the top of the text will be the book and chapter heading in `[book] [chapter]` format, IE "Genesis 15".

### Comparison View

The main pain of the Comparison view will be vertically split, separately scrolling copies of the same book and chapter in separate translations. For example, the left is KJV and the right is ASV. There will be a toggle icon to link both panes to scroll together. Text will be formatted as paragraphs with verse numbers inline.

Similar to the Bible view, at the top and bottom of the center divider will be a "^ Gen 14 ^" to go to the previous chapter, and "v Gen 16 v" to go to the next chapter. The book and chapter heading at the top of each pane will be in `[book] [chapter] ([translation])` format like "Genesis 15 (KJV)".

### Scribe View

The Scribe view has a different purpose. The main pane will consist of three equally sized panes:

1. Bible 1
2. The Scribe Pane
3. Bible 2

Both Bible 1 and Bible 2 shall be in different translations, but the same book and chapter. The book and chapter heading at the top of the text will format like `[book] [chapter] ([translation])` format like "Genesis 15 (KJV)". But the text will be verse-by-verse, with each verse on its own line. The start of a paragraph will be indicated using a ¶ symbol at the start, and will have top margin to separate it from the previous paragraph.

The Scribe Pane will be an editor of the same chapter, but with empty verses that the user fills in with their understanding of the two translations as their own revision in the middle. It will have exactly the same number of verses in the range of both translations. Each verse line will have a slightly darker background color to call it out from the pane when empty. The ¶ symbol and top margin will still be used to indicate a new paragraph. Changes should automatically save after 500ms, and show a brief visual indicator.
