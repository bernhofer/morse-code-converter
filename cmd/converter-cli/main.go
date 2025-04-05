package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bernhofer/morse-code-converter/internal/converter"
)

const (
	morseExtension = ".morse"
	textExtension  = ".txt"
)

type config struct {
	inputFile  string
	outputFile string
}

func main() {
	// Configure logger to not print date/time prefixes for cleaner warnings
	log.SetFlags(0)

	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured when parsing args: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	// Determine conversion direction based on input file extension
	inputExt := strings.ToLower(filepath.Ext(cfg.inputFile))
	isMorseToText := false
	switch inputExt {
	case morseExtension:
		isMorseToText = true
	case textExtension:
		isMorseToText = false
	default:
		fmt.Fprintf(os.Stderr, "Error: Input file must have extension %s or %s\n", textExtension, morseExtension)
		os.Exit(1)
	}

	// Prepare input
	inFile, err := os.Open(cfg.inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file '%s': %v\n", cfg.inputFile, err)
		os.Exit(1)
	}
	defer inFile.Close() // Ensure input file is closed

	// Prepare output
	var outFile io.WriteCloser
	if cfg.outputFile == "" {
		// Use standard output if no file specified
		outFile = os.Stdout
	} else {
		// Use specified output file
		// os.Create truncate existing files
		var err error
		outFile, err = os.Create(cfg.outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating/opening output file '%s': %v\n", cfg.outputFile, err)
			os.Exit(1)
		}

		// Defer closing the output file if not printing to stdout
		defer func() {
			if cerr := outFile.Close(); cerr != nil {
				// Log closing error, as main conversion might have succeeded
				log.Printf("Warning: Error closing output file '%s': %v\n", cfg.outputFile, cerr)
			}
		}()
	}

	// Conversion
	morseConverter := converter.NewMorseConverter(inFile, outFile)
	var convErr error
	if isMorseToText {
		fmt.Println("Conversion from morse to text:")
		convErr = morseConverter.ConvertMorseToText()
	} else {
		fmt.Println("Conversion from test to morse:")
		convErr = morseConverter.ConvertTextToMorse()
	}

	if convErr != nil {
		fmt.Fprintf(os.Stderr, "Error during conversion: %v\n", convErr)
		os.Exit(1)
	}

	fmt.Println("Conversion successful.")
	if cfg.outputFile != "" {
		fmt.Printf("Output written to: %s\n", cfg.outputFile)
	}
}

// Parses command line arguments
func parseFlags() (config, error) {
	var cfg config
	flag.StringVar(&cfg.inputFile, "input", "", "Path to the input file (.txt or .morse) (required)")
	flag.StringVar(&cfg.outputFile, "output", "", "Path to the output file (optional, defaults to stdout)")

	flag.Usage = printUsage // Set custom usage function

	flag.Parse()

	// Validate required flags
	if cfg.inputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Input file path is required.")
		return config{}, fmt.Errorf("missing input file")
	}

	// Check if input exists
	if _, err := os.Stat(cfg.inputFile); os.IsNotExist(err) {
		return config{}, fmt.Errorf("input file '%s' not found", cfg.inputFile)
	} else if err != nil {
		return config{}, fmt.Errorf("error checking input file '%s': %w", cfg.inputFile, err)
	}

	// Check for additional arguments
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Unexpected arguments provided: %v\n", flag.Args())
		return config{}, fmt.Errorf("unexpected arguments")
	}

	return cfg, nil
}

// Print help
func printUsage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s -input <input_file> [-output <output_file>]\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, "\nConverts text (.txt) to Morse code (.morse) or vice versa.")
	fmt.Fprintln(os.Stderr, "The conversion direction is determined by the input file extension.")
	fmt.Fprintln(os.Stderr, "\nArguments:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintf(os.Stderr, "  %s -input message.txt -output message.morse\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s -input code.morse\n", filepath.Base(os.Args[0]))
}
