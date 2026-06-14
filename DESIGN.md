# Boanergies Design

- Our front-end design will be via DaisyUI, using the DaisyUI skill.
- We will introduce each app component as a Svelte component.

## Layouts

- There will always be a left sidebar used to: change views, change translation, change books (in the translation), and change chapter within the current book.
  - The view selector shall be a toggle between single (default / "Bible"), Comparison, and Scribe.
  - Book selection will be a left/right column for Old and New Testaments (ot/nt).
- The default layout will be the Bible view, where the main pane is a single column paragraph Bible reader.
- Two alternate layouts will be: Comparison view, a two column view of the same book/chapter in two different translations; and Scribe view, a three column view with two different translations to the left/right of an editor.
- There should be menu item that allows the user to edit settings, probably to allow control over font-size, line-height, as well as color/contrast.

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