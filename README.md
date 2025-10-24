# name-tabler

**name-tabler** is a small Go utility that generates printable PDF sheets of names
arranged in a grid layout. Each grid cell contains a list of names that can be
cut and used as labels for shelves, boards, or other organized displays.

Each cell is designed to be printed, cut out, and stuck on a shelf or surface
to help people easily find their names.

## Features

- Generates a PDF with configurable layout  
- Adjustable number of columns, rows, and names per cell  
- Automatically paginates if the name list exceeds one page  
- Supports accented Latin characters (é, è, ç, ö, ñ, etc.)  
- Works without any external font dependencies  
- Optional dashed or dotted borders for easy cutting  

## Installation

Using `go install`:

    go install github.com/smallwat3r/name-tabler@latest

Or, this util has no runtime dependencies, you can download a binary for your platform
[here](https://github.com/smallwat3r/name-tabler/releases).

## Usage

1. Prepare a text file named [`names.txt`](./names.txt) in the **current directory**
   where you will run the program. The file must contain one name per line.

2. Run the program to generate a PDF:

        name-tabler

   By default, the PDF will contain:
   - 3 columns
   - 7 rows
   - 5 names per cell

3. To customize the layout:

        name-tabler -cols 2 -rows 5 -n-tab 6

   The generated file will be saved as `names.pdf` in the current directory.

4. Use `-h` to print all options:

        name-tabler -h

   Output:

        Usage of name-tabler:
        -cols int
                number of columns per page (>=1) (default 3)
        -file string
                path to input names file (default "names.txt")
        -n-tab int
                names per table cell (>=1) (default 5)
        -out string
                output PDF file (default "names.pdf")
        -rows int
                number of rows per page (>=1) (default 7)

## Build and Packaging

To build binaries for multiple platforms:

	make build

To package artifacts (`.tar.gz` for Linux/macOS, `.zip` for Windows):

	make pack

All generated files are placed in the `dist/` directory.

## Requirements

- Go 1.20 or later  
- (Optional) GNU Make for building and packaging  

## License

MIT License
