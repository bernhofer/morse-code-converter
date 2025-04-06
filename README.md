
# Morse Code Converter

This Go project is a command-line tool that allows you to convert text to Morse code and Morse code to text. The conversion direction is automatically determined by the extension of the input file.

## Features

* **Text to Morse:** Converts standard text from a `.txt` file to Morse code in a `.morse` file or to the standard output.
* **Morse to Text:** Converts Morse code from a `.morse` file to readable text in a `.txt` file or to the standard output.
* **Command-Line Interface:** Easy to use via the command line.

## Installation

### Prerequisites

* **Go:** You need to have Go installed on your system. [Download Go](https://go.dev/doc/install).

### Steps

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/bernhofer/morse-code-converter.git
    cd morse-code-converter
    ```

2.  **Build the project:**

    ```bash
    go build -o morseconverter ./cmd/converter-cli/main.go
    ```

    This command will compile the `main.go` file and create an executable named `morseconverter` (or `morseconverter.exe` on Windows) in the current directory.

## Usage

The `morseconverter` tool takes an input file path as a required argument and an optional output file path. The conversion direction is determined by the input file's extension.
Non-convertible unicode characters or Morse code sequences will be replaced by the unicode replacement character `�`.

```bash
./morseconverter -input <input_file> [-output <output_file>]
```

### Data Format

This tool expects the input file to adhere to the following format based on its extension:

* **`.txt` files:** These files are expected to contain standard text encoded in UTF-8.
* **`.morse` files:** These files are expected to contain Morse code using the following representation.

**Morse Code Representation:**

When converting to or from Morse code, the following conventions are used:

* **Letters:** Letters within a word are separated by a single space. For example, "hello" would be represented as `.... . .-.. .-.. ---`.
* **Words:** Words are separated by the special character `/`. For example, "hello world" would be represented as `.... . .-.. .-.. --- / .-- --- .-. .-.. -..`.
* **Newlines:** A newline character in the text is encoded as `//` in Morse code.

### Arguments

* `-input <input_file>`: Path to the input file. This file must have either a `.txt` extension for text to Morse conversion or a `.morse` extension for Morse to text conversion. **This argument is required.**
* `-output <output_file>`: (Optional) Path to the output file. If not specified, the output will be printed to the standard output (your terminal). The extension of the output file will be automatically determined based on the conversion direction (e.g., `.morse` for text to Morse, `.txt` for Morse to text). You can also include the correct extension if you like. Existing files will be truncated.

### Examples
Convert Morse code to text and print to stdout:
```bash
./morseconverter -input examples/valid_conv.morse
```

Convert text to Morse code and write it to a new file `output.morse`. The tool will truncate `output.morse` if it exists:
```bash
./morseconverter -input examples/valid_conv.txt  -output ./output
```

