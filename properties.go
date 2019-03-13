package azuretexttospeech

// AudioOutput defines supported audio formats
// Each incorporates a bitrate and encoding type. Azure Speech Service supports 24-KHz, 16-KHz, and 8-KHz audio outputs
// See https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#audio-outputs
type AudioOutput string

const (
	RIFF8Bit8kHzMonoPCM          AudioOutput = "riff-8khz-8bit-mono-mulaw"
	RIFF16Bit16kHzMonoPCM        AudioOutput = "riff-16khz-16bit-mono-pcm"
	RIFF16khz16kbpsMonoSiren     AudioOutput = "riff-16khz-16kbps-mono-siren"
	RIFF24khz16bitMonoPcm        AudioOutput = "riff-24khz-16bit-mono-pcm"
	RAW8Bit8kHzMonoMulaw         AudioOutput = "raw-8khz-8bit-mono-mulaw"
	RAW16Bit16kHzMonoMulaw       AudioOutput = "raw-16khz-16bit-mono-pcm"
	RAW24khz16bitMonoPcm         AudioOutput = "raw-24khz-16bit-mono-pcm"
	Ssml16khz16bitMonoTts        AudioOutput = "ssml-16khz-16bit-mono-tts"
	Audio16khz16kbpsMonoSiren    AudioOutput = "audio-16khz-16kbps-mono-siren"
	Audio16khz32kbitrateMonoMp3  AudioOutput = "audio-16khz-32kbitrate-mono-mp3"
	Audio16khz64kbitrateMonoMp3  AudioOutput = "audio-16khz-64kbitrate-mono-mp3"
	Audio16khz128kbitrateMonoMp3 AudioOutput = "audio-16khz-128kbitrate-mono-mp3"
	Audio24khz48kbitrateMonoMp3  AudioOutput = "audio-24khz-48kbitrate-mono-mp3"
	Audio24khz96kbitrateMonoMp3  AudioOutput = "audio-24khz-96kbitrate-mono-mp3"
)

// Gender type for the digitized language
type Gender string

const (
	// Male , Female are the static Gender constants for digitized voices.
	// See Gender in https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/language-support#standard-voices for breakdown
	Male   Gender = "Male"
	Female Gender = "Female"
)

// Region references the language or locale for text-to-speech.
// See "locale" in https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/language-support#standard-voices
type Region string

const (
	ArEG Region = "ar-EG"
	ArSA Region = "as-SA"
	BgBG Region = "bg-BG"
	CaES Region = "ca-ES"
	CsCZ Region = "cs-CZ"
	DaDK Region = "da-DK"
	DeAT Region = "de-AT"
	DeCH Region = "de-CH"
	DeDE Region = "de-DE"
	ElGR Region = "el-GR"
	EnAU Region = "en-AU"
	EnCA Region = "en-CA"
	EnGB Region = "en-GB"
	EnIE Region = "en-IE"
	EnIN Region = "en-IN"
	EnUS Region = "en-US"
	EsES Region = "es-ES"
	EsMX Region = "es-MX"
	FiFI Region = "fi-FI"
	FrCA Region = "fr-CA"
	FrCH Region = "fr-CH"
	FrFR Region = "fr-FR"
	HeIL Region = "he-IL"
	HiIN Region = "hi-IN"
	HrHR Region = "hr-HR"
	HuHU Region = "hu-HU"
	IdID Region = "id-ID"
	ItIT Region = "it-IT"
	JaJP Region = "ja-JP"
	KoKR Region = "ko-KR"
	MsMY Region = "ms-MY"
	NbNO Region = "nb-NO"
	NlNL Region = "nl-NL"
	PlPL Region = "pl-PL"
	PtBR Region = "pt-BR"
	PtPT Region = "pt-PT"
	RoRO Region = "ro-RO"
	RuRU Region = "ru-RU"
	SkSK Region = "sk-SK"
	SlSL Region = "sl-SL"
	SvSE Region = "sv-SE"
	TaIN Region = "ta-IN"
	TeIN Region = "te-IN"
	ThTH Region = "th-TH"
	TrTR Region = "tr-TR"
	ViVN Region = "vi-VN"
	ZhCN Region = "zh-CN"
	ZhHK Region = "zh-HK"
	ZhTW Region = "zh-TW"
)

// localeGender represents the key used within the `localeToGender` map.
type localeGender struct {
	locale Region
	gender Gender
}

// localeToGender maps a language|locale and voice gender to the description string
var localeToGender = map[localeGender]string{
	localeGender{ArEG, Female}: "(ar-EG, Hoda)",
	localeGender{ArSA, Male}:   "(ar-SA, Naayf)",
	localeGender{BgBG, Male}:   "(bg-BG, Ivan)",
	localeGender{CaES, Female}: "(ca-ES, HerenaRUS)",
	localeGender{CsCZ, Male}:   "(cs-CZ, Jakub)",
	localeGender{DaDK, Female}: "(da-DK, HelleRUS)",
	localeGender{DeAT, Male}:   "(de-AT, Michael)",
	localeGender{DeCH, Male}:   "(de-CH, Karsten)",
	localeGender{DeDE, Female}: "(de-DE, Hedda)",
	localeGender{DeDE, Male}:   "(de-DE, Stefan, Apollo)",
	localeGender{ElGR, Male}:   "(el-GR, Stefanos)",
	localeGender{EnAU, Female}: "(en-AU, Catherine)",
	localeGender{EnCA, Female}: "(en-CA, Linda)",
	localeGender{EnGB, Female}: "(en-GB, Susan, Apollo)",
	localeGender{EnGB, Male}:   "(en-GB, George, Apollo)",
	localeGender{EnIE, Male}:   "(en-IE, Sean)",
	localeGender{EnIN, Female}: "(en-IN, Heera, Apollo)",
	localeGender{EnIN, Male}:   "(en-IN, Ravi, Apollo)",
	localeGender{EnUS, Female}: "(en-US, ZiraRUS)",
	localeGender{EnUS, Male}:   "(en-US, BenjaminRUS)",
	localeGender{EsES, Female}: "(es-ES, Laura, Apollo)",
	localeGender{EsES, Male}:   "(es-ES, Pablo, Apollo)",
	localeGender{EsMX, Female}: "(es-MX, HildaRUS)",
	localeGender{EsMX, Male}:   "(es-MX, Raul, Apollo)",
	localeGender{FiFI, Female}: "(fi-FI, HeidiRUS)",
	localeGender{FrCA, Female}: "(fr-CA, Caroline)",
	localeGender{FrCH, Male}:   "(fr-CH, Guillaume)",
	localeGender{FrFR, Female}: "(fr-FR, Julie, Apollo)",
	localeGender{FrFR, Male}:   "(fr-FR, Paul, Apollo)",
	localeGender{HeIL, Male}:   "(he-IL, Asaf)",
	localeGender{HiIN, Female}: "(hi-IN, Kalpana, Apollo)",
	localeGender{HiIN, Male}:   "(hi-IN, Hemant)",
	localeGender{HrHR, Male}:   "(hr-HR, Matej)",
	localeGender{HuHU, Male}:   "(hu-HU, Szabolcs)",
	localeGender{IdID, Male}:   "(id-ID, Andika)",
	localeGender{ItIT, Male}:   "(it-IT, Cosimo, Apollo)",
	localeGender{ItIT, Female}: "(it-IT, LuciaRUS)",
	localeGender{JaJP, Male}:   "(ja-JP, Ichiro, Apollo)",
	localeGender{JaJP, Female}: "(ja-JP, HarukaRUS)",
	localeGender{KoKR, Female}: "(ko-KR, HeamiRUS)",
	localeGender{MsMY, Male}:   "(ms-MY, Rizwan)",
	localeGender{NbNO, Female}: "nb-NO, HuldaRUS)",
	localeGender{NlNL, Female}: "nl-NL, HannaRUS)",
	localeGender{PlPL, Female}: "(pl-PL, PaulinaRUS)",
	localeGender{PtBR, Female}: "pt-BR, HeloisaRUS)",
	localeGender{PtBR, Male}:   "(pt-BR, Daniel, Apollo)",
	localeGender{PtPT, Female}: "(pt-PT, HeliaRUS)",
	localeGender{RoRO, Male}:   "(ro-RO, Andrei)",
	localeGender{RuRU, Female}: "ru-RU, Irina, Apollo)",
	localeGender{RuRU, Male}:   "(ru-RU, Pavel, Apollo)",
	localeGender{SkSK, Male}:   "(sk-SK, Filip)",
	localeGender{SlSL, Male}:   "(sl-SI, Lado)",
	localeGender{SvSE, Female}: "(sv-SE, HedvigRUS)",
	localeGender{TaIN, Male}:   "(ta-IN, Valluvar)",
	localeGender{TeIN, Female}: "(te-IN, Chitra)",
	localeGender{ThTH, Male}:   "th-TH, Pattara)",
	localeGender{TrTR, Female}: "(tr-TR, SedaRUS)",
	localeGender{ViVN, Male}:   "(vi-VN, An)",
	localeGender{ZhCN, Female}: "(zh-CN, HuihuiRUS)",
	localeGender{ZhCN, Male}:   "(zh-CN, Kangkang, Apollo)",
	localeGender{ZhHK, Female}: "(zh-HK, Tracy, Apollo)",
	localeGender{ZhHK, Male}:   "(zh-HK, Danny, Apollo)",
	localeGender{ZhTW, Female}: "(zh-TW, Yating, Apollo)",
	localeGender{ZhTW, Male}:   "(zh-TW, Zhiwei, Apollo)",
}
