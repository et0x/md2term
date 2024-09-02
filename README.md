# md2term

md2term is a command-line tool that renders Markdown content in the terminal with syntax highlighting and formatting.

## Features

- Supports common Markdown elements (headers, lists, code blocks, tables, etc.)
- Colorful and styled output in the terminal
- Can read from files or stdin
- Easy to use and integrate into your workflow

## Installation

To install md2term, make sure you have Go installed on your system, then run:

```
go install github.com/et0x/md2term/cmd/md2term@latest
```

## Usage

md2term reads from stdin by default, but you can specify a file to read from.

```
echo "# Hello, World!" | md2term

cat README.md | md2term
```

If you don't specify a file, md2term will read from stdin.

## License

md2term is licensed under the MIT License. See the LICENSE file for more details.
