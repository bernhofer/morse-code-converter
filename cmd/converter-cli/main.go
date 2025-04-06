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

// Valid file extensions
const (
	morseExtension = ".morse"
	textExtension  = ".txt"
)

// User userConfig for cli args
type userConfig struct {
	inputFile  string
	outputFile string
}

func main() {
	// Configure logger to not print date/time prefixes for cleaner warnings
	log.SetFlags(0)

	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured when parsing args: %v\n", err)
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
		// Check output file extension
		var outFileBasename string = filepath.Base(cfg.outputFile)
		var validOutExtension bool = strings.HasSuffix(outFileBasename, ".txt") || strings.HasSuffix(outFileBasename, ".morse")
		if !validOutExtension {
			if strings.Contains(outFileBasename, ".") {
				// Some extension other than txt or morse is already defined
				fmt.Fprintf(os.Stderr, "Error: invalid or missing extension for output file '%s'\n", cfg.outputFile)
				os.Exit(1)
			} else if isMorseToText {
				// Set output extension to inverse of input extension
				cfg.outputFile = cfg.outputFile + ".txt"
			} else {
				cfg.outputFile = cfg.outputFile + ".morse"
			}
		}
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
func parseFlags() (userConfig, error) {
	var cfg userConfig
	flag.StringVar(&cfg.inputFile, "input", "", "Path to the input file (.txt or .morse) (required)")
	flag.StringVar(&cfg.inputFile, "i", "", "Alias for -input")

	flag.StringVar(&cfg.outputFile, "output", "", "Path to the output file (optional, defaults to stdout)")
	flag.StringVar(&cfg.outputFile, "o", "", "Alias for -output")

	// Set custom printUsage() function
	flag.Usage = printUsage

	flag.Parse()

	// Validate required flags
	if cfg.inputFile == "" {
		printUsage()
		return userConfig{}, fmt.Errorf("missing input file")
	}

	// Check if input file exists
	if _, err := os.Stat(cfg.inputFile); os.IsNotExist(err) {
		return userConfig{}, fmt.Errorf("input file '%s' not found", cfg.inputFile)
	} else if err != nil {
		return userConfig{}, fmt.Errorf("error checking input file '%s': %w", cfg.inputFile, err)
	}

	// Check for additional positional flags
	// Undefined flags (e.g. -a xyz) are already captured by flag.Parse()
	if len(flag.Args()) > 0 {
		printUsage()
		return userConfig{}, fmt.Errorf("unexpected arguments")
	}

	return cfg, nil
}

// Print help
func printUsage() {
	fmt.Fprintln(os.Stderr, "\n############## morse-code-converter help ###############")
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, "\nConverts text (.txt) to Morse code (.morse) or vice versa.")
	fmt.Fprintln(os.Stderr, "The conversion direction is determined by the input file extension.")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintf(os.Stderr, "  %s -input message.txt -output message.morse\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s -i code.morse\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, "########################################################")
}
