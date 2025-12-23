# Samla

A user-friendly desktop app for organizing and managing your collections. Perfect for crafting supplies, stamps, dies, tools, or any items stored in boxes and bags!

**Samla** (Swedish for "to collect") makes it easy to catalog, find, and organize your sets.

## Features

- **Fuzzy Search** – Find sets by name, tags, products, manufacturer, box, or bag number using intelligent fuzzy matching (powered by Fuse.js)
- **Special Filters** – Use `@Box`, `@Product`, `@Manufacturer`, `@Tag`, or `@Location` prefixes for targeted searches
- **Storage Management** – Organize items by locations, boxes, and bags with room/shelf/compartment details
- **Product Tracking** – Record individual items within each set with type classification
- **Tags & Categories** – Organize with keywords, manufacturers, and set types
- **Images** – Add photos for quick visual identification
- **Overview Mode** – Click on a set to see a beautiful overview before editing
- **Sorting Options** – Sort by name, box, location, or newest first
- **User-Friendly** – Clean interface with large text and intuitive navigation

## Data Storage

- Base folder: `%APPDATA%/Samla` (Windows) or `~/.config/Samla` (Linux/macOS)
  - `Data/samla.db` – SQLite database
  - `Images/` – Stored images
- Open the folder directly from the app using the folder icon in the header.

## Search Examples

| Query              | Description                                        |
| ------------------ | -------------------------------------------------- |
| `flower`           | Fuzzy search across all fields                     |
| `@Box A01`         | Find all sets in box A01                           |
| `@Product rose`    | Find sets containing a product named "rose"        |
| `@Manufacturer CP` | Find sets by manufacturer                          |
| `@Tag christmas`   | Find sets tagged with "christmas"                  |
| `@Location office` | Find sets stored in a location containing "office" |

## Keyboard Shortcuts

- `Ctrl+F` – Focus search bar
- `Ctrl+N` – Create new set
- `Escape` – Go back to previous view
- `Delete` – Remove selected product

## Development

Built with:

- **Backend**: Go 1.24+, Wails 2.11, modernc.org/sqlite
- **Frontend**: Vue 3, Vite, TypeScript, Fuse.js

### Live Development

```bash
wails dev
```

### Build

```bash
wails build
```

## License

MIT
