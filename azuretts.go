package azuretexttospeech

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// v1Path is the URL path *v1* of Azure's cognitive services
const v1Path = "/cognitiveservices/v1"

// digitizeActionTimeout is the amount of time the http client will wait for a response during digitize request
const digitizeActionTimeout = time.Second * 30

// refreshTokenTimeout is the amount of time the http client will wait during the token refresh action.
const refreshTokenTimeout = time.Second * 15

// AzureCognitiveServicesAPI references the locations of the availability of standard voices.
// See https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/regions#standard-voices
type AzureCognitiveServicesAPI string

const (
	// Azure regions and their endpoints that support the Text To Speech service.
	WestUS        AzureCognitiveServicesAPI = "https://westus.tts.speech.microsoft.com"
	WestUS2       AzureCognitiveServicesAPI = "https://westus2.tts.speech.microsoft.com"
	EastUS        AzureCognitiveServicesAPI = "https://eastus.tts.speech.microsoft.com"
	EastUS2       AzureCognitiveServicesAPI = "https://eastus2.tts.speech.microsoft.com"
	EastAsia      AzureCognitiveServicesAPI = "https://eastasia.tts.speech.microsoft.com"
	SoutheastAsia AzureCognitiveServicesAPI = "https://southeastasia.tts.speech.microsoft.com"
	NorthEurpoe   AzureCognitiveServicesAPI = "https://northeurope.tts.speech.microsoft.com"
	WestEurope    AzureCognitiveServicesAPI = "https://westeurope.tts.speech.microsoft.com"
)

// AzureCognitiveServicesToken references the Azure endpoints required to fetch the bearer token.
// https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#how-to-get-an-access-token
type AzureCognitiveServicesToken string

const (
	// 	Endpoints for token refresh/assignment.
	WestUSToken        AzureCognitiveServicesToken = "https://westus.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	WestUS2Token       AzureCognitiveServicesToken = "https://westus2.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	EastUSToken        AzureCognitiveServicesToken = "https://eastus.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	EastUS2Token       AzureCognitiveServicesToken = "https://eastus2.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	EastAsiaToken      AzureCognitiveServicesToken = "https://eastasia.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	SoutheastAsiaToken AzureCognitiveServicesToken = "https://southeastasia.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	NorthEuropeToken   AzureCognitiveServicesToken = "https://northeurope.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	WestEuropeToken    AzureCognitiveServicesToken = "https://westeurope.api.cognitive.microsoft.com/sts/v1.0/issueToken"
)

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
)

// localeGender represents the key used within the `localeToGender` map.
type localeGender struct {
	locale Region
	gender Gender
}

// localeToGender maps a language|locale and voice gender to the description string
var localeToGender = map[localeGender]string{
	localeGender{ArEG, Female}: "Microsoft Server Speech Text to Speech Voice (ar-EG, Hoda)",
	localeGender{ArSA, Male}:   "Microsoft Server Speech Text to Speech Voice (ar-SA, Naayf)",
	localeGender{BgBG, Male}:   "Microsoft Server Speech Text to Speech Voice (bg-BG, Ivan)",
	localeGender{CaES, Female}: "Microsoft Server Speech Text to Speech Voice (ca-ES, HerenaRUS)",
	localeGender{CsCZ, Male}:   "Microsoft Server Speech Text to Speech Voice (cs-CZ, Jakub)",
	localeGender{DaDK, Female}: "Microsoft Server Speech Text to Speech Voice (da-DK, HelleRUS)",
	localeGender{DeAT, Male}:   "Microsoft Server Speech Text to Speech Voice (de-AT, Michael)",
	localeGender{DeCH, Male}:   "Microsoft Server Speech Text to Speech Voice (de-CH, Karsten)",
	localeGender{DeDE, Female}: "Microsoft Server Speech Text to Speech Voice (de-DE, Hedda)",
	localeGender{DeDE, Male}:   "Microsoft Server Speech Text to Speech Voice (de-DE, Stefan, Apollo)",
	localeGender{ElGR, Male}:   "Microsoft Server Speech Text to Speech Voice (el-GR, Stefanos)",
	localeGender{EnAU, Female}: "Microsoft Server Speech Text to Speech Voice (en-AU, Catherine)",
	localeGender{EnCA, Female}: "Microsoft Server Speech Text to Speech Voice (en-CA, Linda)",
	localeGender{EnGB, Female}: "Microsoft Server Speech Text to Speech Voice (en-GB, Susan, Apollo)",
	localeGender{EnGB, Male}:   "Microsoft Server Speech Text to Speech Voice (en-GB, George, Apollo)",
	localeGender{EnIE, Male}:   "Microsoft Server Speech Text to Speech Voice (en-IE, Sean)",
	localeGender{EnIN, Female}: "Microsoft Server Speech Text to Speech Voice (en-IN, Heera, Apollo)",
	localeGender{EnIN, Male}:   "Microsoft Server Speech Text to Speech Voice (en-IN, Ravi, Apollo)",
	localeGender{EnUS, Female}: "Microsoft Server Speech Text to Speech Voice (en-US, ZiraRUS)",
	localeGender{EnUS, Male}:   "Microsoft Server Speech Text to Speech Voice (en-US, BenjaminRUS)",
}

// voiceXML renders the XML payload for the TTS api.
// For API reference see https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#sample-request-1
func voiceXML(speechText, description string, locale Region, gender Gender) string {

	return fmt.Sprintf(`<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'>%s</voice></speak>`,
		locale, locale, gender, description, speechText)
}

// Digitize places a call to the TTS api and fetches the audio bytestream. `speechText` is the target string to be
// converted to a digitized recording.
func (az *AzureSpeech) Digitize(speechText string, region Region, gender Gender, audioOutput AudioOutput) ([]byte, error) {

	descriprtion, ok := localeToGender[localeGender{region, gender}]
	if !ok {
		return nil, fmt.Errorf("unable to locale region=%s, gender=%s pair", region, gender)
	}

	v := voiceXML(speechText, descriprtion, region, gender)
	request, err := http.NewRequest(http.MethodPost, az.APIBase, bytes.NewBufferString(v))
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Microsoft-OutputFormat", string(audioOutput))
	request.Header.Set("Content-Type", "application/ssml+xml")
	request.Header.Set("Authorization", "Bearer "+az.accessToken)
	request.Header.Set("User-Agent", "azuretts")

	client := &http.Client{Timeout: digitizeActionTimeout}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code; status=%s", response.Status)
	}
	return ioutil.ReadAll(response.Body)
}

// RefreshToken fetches an updated bearer token from the Azure cognitive speech/text services, or an error if unable to retrive.
// Each token is valid for a maximum of 10 minutes. Details for auth tokens are referenced at
// https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#authentication .
// This does not need to be called
func (az *AzureSpeech) RefreshToken() error {
	request, _ := http.NewRequest(http.MethodPost, az.tokenURL, nil)
	request.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	client := &http.Client{Timeout: refreshTokenTimeout}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code; received status=%s", response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)
	az.accessToken = string(body)
	return nil
}

// startRefreshWorker refreshes the authentication token on at a 9 minute interval. A channel is returned
// if the caller wishes to cancel the channel.
func (az *AzureSpeech) startRefreshWorker() chan bool {
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(time.Minute * 9)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := az.RefreshToken()
				if err != nil {
					log.Printf("failed to refresh token, err=%+v", err)
				}
			case <-done:
				return
			}
		}
	}()
	return done
}

// AzureSpeech stores configuration and state information for the TTS client.
type AzureSpeech struct {
	accessToken     string    // Token received from `AzureCognitiveServicesToken`
	APIBase         string    // target `AzureCognitiveServicesAPI`
	SubscriptionKey string    // API key for Azure's Congnitive Speech services
	TokenRefresher  chan bool // channel to stop the token refresh goroutine.
	tokenURL        string    // Local endpoints for AzureCognitiveServicesToken

}

// New returns an `AzureSpeech` object and starts a background token refresh timer
func New(subscriptionKey string, api AzureCognitiveServicesAPI, token AzureCognitiveServicesToken) (*AzureSpeech, error) {

	az := &AzureSpeech{
		APIBase:         string(api) + v1Path,
		SubscriptionKey: subscriptionKey,
		tokenURL:        string(token)}

	// api requires that the token is refreshed every 10 mintutes.
	// We will do this task in the background every ~9 minutes.
	if err := az.RefreshToken(); err != nil {
		return nil, fmt.Errorf("failed to fetch initial token, err='%+v'", err)
	}
	az.TokenRefresher = az.startRefreshWorker()
	return az, nil
}
