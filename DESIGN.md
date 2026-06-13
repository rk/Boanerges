# Boanergies Design

- Our front-end design will be via DaisyUI, using the DaisyUI skill.
- We will introduce each app component as a Svelte component.

## Layouts

- There will always be a left sidebar used to: change views, change translation, change books (in the translation), and change chapter within the current book.
  - The view selector shall be a toggle between single (default / "Bible"), Comparison, and Scribe.
  - Book selection will be a left/right column for Old and New Testaments (ot/nt).
- The default layout will be the Bible view, where the main pane is a single column paragraph Bible reader.
- Two alternate layouts will be: Comparison view, a two column view of the same book/chapter in two different translations; and Scribe view, a three column view with two different translations to the left/right of an editor.

### Bible View