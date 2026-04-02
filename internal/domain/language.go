package domain

// TargetLanguage defines the supported target languages for practice.
// These are the only languages users can select as their learning target.
type TargetLanguage string

const (
	LangEnglish  TargetLanguage = "english"
	LangFrench   TargetLanguage = "french"
	LangSpanish  TargetLanguage = "spanish"
	LangRussian  TargetLanguage = "russian"
	LangMandarin TargetLanguage = "mandarin"
	LangArabic   TargetLanguage = "arabic"
)

// ValidTargetLanguages is the authoritative list of supported target languages.
var ValidTargetLanguages = []TargetLanguage{
	LangEnglish,
	LangFrench,
	LangSpanish,
	LangRussian,
	LangMandarin,
	LangArabic,
}

// IsValidTargetLanguage checks whether a given string is one of the allowed target languages.
func IsValidTargetLanguage(lang string) bool {
	for _, valid := range ValidTargetLanguages {
		if string(valid) == lang {
			return true
		}
	}
	return false
}

// TargetLanguageNames returns human-readable display names (Indonesian labels) for each language.
func TargetLanguageNames() map[TargetLanguage]string {
	return map[TargetLanguage]string{
		LangEnglish:  "Inggris",
		LangFrench:   "Prancis",
		LangSpanish:  "Spanyol",
		LangRussian:  "Rusia",
		LangMandarin: "Tionghoa (Mandarin)",
		LangArabic:   "Arab",
	}
}
