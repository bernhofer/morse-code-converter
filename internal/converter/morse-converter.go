package converter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"unicode"

	"github.com/bernhofer/morse-code-converter/internal/mappings"
)

// Replace unknown/unexpected chars in text/morse with replacement character
const usePlaceholderForTextOutput bool = true   // Enable because txt is utf-8
const usePlaceholderForMorseOutput bool = false // Disabled to ensure valid morse code

// Struct that handles the conversion
type MorseConverter struct {
	Input  io.Reader // Data source (morse or txt file)
	Output io.Writer // Destination for converted data (txt or morse file)
}

// Instantiation function for converter
func NewMorseConverter(input io.Reader, output io.Writer) *MorseConverter {
	return &MorseConverter{
		Input:  input,
		Output: output,
	}
}

// ConvertTextToMorse reads UTF-8 text from the input and writes Morse code to the output.
func (c *MorseConverter) ConvertTextToMorse() error {
	// Create reader & writer
	reader := bufio.NewReader(c.Input)
	writer := bufio.NewWriter(c.Output)

	// Ensure buffer is flushed at the end
	defer writer.Flush()

	// Track beginning of char sequence to avoid leading space
	firstChar := true

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break // End of file → done with processing
			}
			return fmt.Errorf("error reading input: %w", err)
		}

		// Ignore carriage return characters
		if r == '\r' {
			continue
		}

		// Check if character is in mapping, mapping only contains capitalized letters
		upperRune := unicode.ToUpper(r)
		morseCode, validMapping := mappings.TextToMorse[upperRune]

		if validMapping {
			// Found mapping
			if !firstChar {
				// Add space separator before the next code, unless it's the very first one
				_, err = writer.WriteString(" ")
				if err != nil {
					return fmt.Errorf("error writing separator: %w", err)
				}
			}
			_, err = writer.WriteString(morseCode)
			if err != nil {
				return fmt.Errorf("error writing morse code for '%c': %w", r, err)
			}
			firstChar = false
		} else {
			// Handle unsupported characters
			log.Printf("Warning: skipping unsupported character: '%c' (U+%04X)\n", r, r)
			// Optionally write a placeholder - uncomment if desired
			if usePlaceholderForMorseOutput {
				if !firstChar {
					_, err = writer.WriteString(" ")
					if err != nil {
						return fmt.Errorf("error writing separator for placeholder: %w", err)
					}
				}
				_, err = writer.WriteRune(mappings.UnknownMorse) // e.g., "?"
				if err != nil {
					return fmt.Errorf("error writing placeholder morse: %w", err)
				}
				firstChar = false
			}
		}
	}

	// Add a final newline to the output for cleaner formatting in terminals/files
	_, err := writer.WriteString("\n")
	if err != nil {
		// Log error, but don't fail the whole process just for the final newline
		log.Printf("Warning: could not write final newline: %v", err)
	}

	return nil
}

// ConvertMorseToText reads Morse code from the input and writes text to the output.
func (c *MorseConverter) ConvertMorseToText() error {
	// Use scanner to read space-separated morse sequence
	scanner := bufio.NewScanner(c.Input)
	scanner.Split(bufio.ScanWords)

	// Create writer
	writer := bufio.NewWriter(c.Output)
	defer writer.Flush() // Ensure buffer is flushed at the end

	for scanner.Scan() {
		morseWord := scanner.Text()

		// Handle special Morse codes first
		switch morseWord {
		case "/": // Word separator
			_, err := writer.WriteString(" ") // Output a single space
			if err != nil {
				return fmt.Errorf("error writing space for '/': %w", err)
			}
		case "//": // Newline
			_, err := writer.WriteString("\n") // Output a newline
			if err != nil {
				return fmt.Errorf("error writing newline for '//': %w", err)
			}
		default:
			// Look up the standard Morse code
			if textChar, ok := mappings.MorseToText[morseWord]; ok {
				_, err := writer.WriteRune(textChar)
				if err != nil {
					return fmt.Errorf("error writing character for '%s': %w", morseWord, err)
				}
			} else {
				if usePlaceholderForTextOutput {
					// Handle unrecognized Morse sequence
					log.Printf("Warning: Replacing unrecognized morse sequence: '%s'\n", morseWord)
					_, err := writer.WriteRune(mappings.UnknownChar) // Write '?'
					if err != nil {
						return fmt.Errorf("error writing placeholder char: %w", err)
					}
				} else {
					log.Printf("Warning: Skipping unrecognized morse sequence: '%s'\n", morseWord)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning input: %w", err)
	}

	// Add a final newline
	_, err := writer.WriteString("\n")
	if err != nil {
		log.Printf("Warning: could not write final newline: %v", err)
	}

	return nil
}
