package azuretexttospeech

// AudioOutput defines supported audio formats
// Each incorporates a bitrate and encoding type. Azure Speech Service supports 24-KHz, 16-KHz, and 8-KHz audio outputs
// See https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#audio-outputs
type AudioOutput int

const (
	AudioRIFF8Bit8kHzMonoPCM AudioOutput = iota
	AudioRIFF16Bit16kHzMonoPCM
	AudioRIFF16khz16kbpsMonoSiren
	AudioRIFF24khz16bitMonoPcm
	AudioRAW8Bit8kHzMonoMulaw
	AudioRAW16Bit16kHzMonoMulaw
	AudioRAW24khz16bitMonoPcm
	AudioSsml16khz16bitMonoTts
	Audio16khz16kbpsMonoSiren
	Audio16khz32kbitrateMonoMp3
	Audio6khz64kbitrateMonoMp3
	Audio16khz128kbitrateMonoMp3
	Audio24khz48kbitrateMonoMp3
	Audio24khz96kbitrateMonoMp3
)

func (a AudioOutput) String() string {
	return []string{"riff-8khz-8bit-mono-mulaw",
		"riff-16khz-16bit-mono-pcm",
		"riff-16khz-16kbps-mono-siren",
		"riff-24khz-16bit-mono-pcm",
		"raw-8khz-8bit-mono-mulaw",
		"raw-16khz-16bit-mono-pcm",
		"raw-24khz-16bit-mono-pcm",
		"ssml-16khz-16bit-mono-tts",
		"audio-16khz-16kbps-mono-siren",
		"audio-16khz-32kbitrate-mono-mp3",
		"audio-16khz-64kbitrate-mono-mp3",
		"audio-16khz-128kbitrate-mono-mp3",
		"audio-24khz-48kbitrate-mono-mp3",
		"audio-24khz-96kbitrate-mono-mp3",
	}[a]
}

// Gender type for the digitized language
//go:generate enumer -type=Gender -linecomment -json
type Gender int

const (
	// GenderMale , GenderFemale are the static Gender constants for digitized voices.
	// See Gender in https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/language-support#standard-voices for breakdown
	GenderMale   Gender = iota // Male
	GenderFemale               // Female
)

// Locale references the language or locale for text-to-speech.
// See "locale" in https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/language-support#standard-voices
//go:generate enumer -type=Locale -linecomment -json
type Locale int

const (
	LocaleArEG Locale = iota //ar-EG
	LocaleArSA               // ar-SA
	LocaleBgBG               // bg-BG
	LocaleCaES               // ca-ES
	LocaleCsCZ               // cs-CZ
	LocaleDaDK               // da-DK
	LocaleDeAT               // de-AT
	LocaleDeCH               // de-CH
	LocaleDeDE               // de-DE
	LocaleElGR               // el-GR
	LocaleEnAU               // en-AU
	LocaleEnCA               // en-CA
	LocaleEnGB               // en-GB
	LocaleEnIE               // en-IE
	LocaleEnIN               // en-IN
	LocaleEnUS               // en-US
	LocaleEsES               // es-ES
	LocaleEsMX               // es-MX
	LocaleEtEE               // et-EE
	LocaleFiFI               // fi-FI
	LocaleFrCA               // fr-CA
	LocaleFrCH               // fr-CH
	LocaleFrFR               // fr-FR
	LocaleGaIE               // ga-IE
	LocaleHeIL               // he-IL
	LocaleHiIN               // hi-IN
	LocaleHrHR               // hr-HR
	LocaleHuHU               // hu-HU
	LocaleIdID               // id-ID
	LocaleItIT               // it-IT
	LocaleJaJP               // ja-JP
	LocaleKoKR               // ko-KR
	LocaleLtLT               // lt-LT
	LocaleLvLV               // lv-LV
	LocaleMtMT               // mt-MT
	LocaleMrIN               // mr-IN
	LocaleMsMY               // ms-MY
	LocaleNbNO               // nb-NO
	LocaleNlNL               // nl-NL
	LocalePlPL               // pl-PL
	LocalePtBR               // pt-BR
	LocalePtPT               // pt-PT
	LocaleRoRO               // ro-RO
	LocaleRuRU               // ru-RU
	LocaleSkSK               // sk-SK
	LocaleSlSI               // sl-SI
	LocaleSvSE               // sv-SE
	LocaleTaIN               // ta-IN
	LocaleTeIN               // te-IN
	LocaleThTH               // th-TH
	LocaleTrTR               // tr-TR
	LocaleViVN               // vi-VN
	LocaleZhCN               // zh-CN
	LocaleZhHK               // zh-HK
	LocaleZhTW               // zh-TW
)
