# md2term

md2term is a command-line tool that renders Markdown content in the terminal with syntax highlighting and formatting.

## Features

- Supports common Markdown elements (headers, lists, code blocks, tables, etc.)
- Colorful and styled output in the terminal
- Can read from files or stdin
- Easy to use and integrate into your workflow
- Customizable themes for different visual preferences

## Installation

To install md2term, make sure you have Go installed on your system, then run:

```
go install github.com/et0x/md2term/cmd/md2term@latest
```

## Usage

md2term reads from stdin by default, but you can specify a file to read from:

```
echo "# Hello, World!" | md2term

cat README.md | md2term

md2term README.md
```

### Theme Selection

md2term supports multiple themes for different visual preferences. You can specify a theme using the `-theme` flag:

```
md2term -theme dark README.md
```

Available themes:

- `default`: A colorful theme suitable for most terminal backgrounds
- `dark`: A theme optimized for dark terminal backgrounds
- `light`: A theme optimized for light terminal backgrounds

### Setting a Default Theme

You can set a theme as your default using the `-set-default-theme` flag:

```
md2term -theme dark -set-default-theme
```

This will save your theme preference for future use.

### Version Information

To check the version of md2term, use the `-version` flag:

```
md2term -version
```

## Customizing Themes

md2term uses a `themes.yaml` file to define color schemes. You can modify this file to create your own custom themes. The theme structure includes colors for various Markdown elements:

- Headings (levels 1-6)
- List items (regular, checked, unchecked)
- Table elements (border, cells, header, even/odd rows)
- Code blocks and inline code
- Links and images
- Blockquotes

Refer to the `themes.yaml` file in the repository for the complete theme structure and color definitions.

## License

md2term is licensed under the MIT License. See the LICENSE file for more details.
