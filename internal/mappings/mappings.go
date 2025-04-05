package mappings

// Based on International Morse Code standard
// https://en.wikipedia.org/wiki/Morse_code#Letters,_numbers,_punctuation,_prosigns_for_Morse_code_and_non-Latin_variants
// Public map
var TextToMorse = map[rune]string{
	// General Morse code
	// Letters
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.",
	'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.",
	'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",
	// Accented letters
	'È': "..-..",
	// Numbers
	'1': ".----", '2': "..---", '3': "...--", '4': "....-", '5': ".....",
	'6': "-....", '7': "--...", '8': "---..", '9': "----.", '0': "-----",
	// Punctuation
	'.': ".-.-.-", ',': "--..--", '?': "..--..", '\'': ".----.", '!': "-.-.--",
	'/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...", ':': "---...",
	';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-", '_': "..--.-",
	'"': ".-..-.", '$': "...-..-", '@': ".--.-.",
	// Special characters for this task
	' ':  "/",  // Word separator
	'\n': "//", // Newline
}

// MorseToText is generated automatically from TextToMorse for consistency
var MorseToText = map[string]rune{}

// Placeholder for characters that cannot be converted
const UnknownChar rune = '\uFFFD'
const UnknownMorse rune = '\uFFFD'

// init generates the reverse map (morseToText) from TextToMorse
func init() {
	for r, s := range TextToMorse {
		// Ensure we don't overwrite existing entries if Morse codes were accidentally duplicated
		// (though the standard map shouldn't have duplicates)
		if _, exists := MorseToText[s]; !exists {
			MorseToText[s] = r
		}
	}
}
