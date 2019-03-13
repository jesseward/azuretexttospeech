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

// synthesizeActionTimeout is the amount of time the http client will wait for a response during Synthesize request
const synthesizeActionTimeout = time.Second * 30

// tokenRefreshTimeout is the amount of time the http client will wait during the token refresh action.
const tokenRefreshTimeout = time.Second * 15

// TextToSpeechAPI references the locations of the availability of standard voices.
// See https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/regions#standard-voices
type TextToSpeechAPI string

const (
	// Azure regions and their endpoints that support the Text To Speech service.
	WestUS        TextToSpeechAPI = "https://westus.tts.speech.microsoft.com"
	WestUS2       TextToSpeechAPI = "https://westus2.tts.speech.microsoft.com"
	EastUS        TextToSpeechAPI = "https://eastus.tts.speech.microsoft.com"
	EastUS2       TextToSpeechAPI = "https://eastus2.tts.speech.microsoft.com"
	EastAsia      TextToSpeechAPI = "https://eastasia.tts.speech.microsoft.com"
	SoutheastAsia TextToSpeechAPI = "https://southeastasia.tts.speech.microsoft.com"
	NorthEurpoe   TextToSpeechAPI = "https://northeurope.tts.speech.microsoft.com"
	WestEurope    TextToSpeechAPI = "https://westeurope.tts.speech.microsoft.com"
)

// TokenRefreshAPI references the Azure endpoints required to fetch the bearer token.
// https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#how-to-get-an-access-token
type TokenRefreshAPI string

const (
	// 	Endpoints for token refresh/assignment.
	WestUSToken        TokenRefreshAPI = "https://westus.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	WestUS2Token       TokenRefreshAPI = "https://westus2.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	EastUSToken        TokenRefreshAPI = "https://eastus.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	EastUS2Token       TokenRefreshAPI = "https://eastus2.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	EastAsiaToken      TokenRefreshAPI = "https://eastasia.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	SoutheastAsiaToken TokenRefreshAPI = "https://southeastasia.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	NorthEuropeToken   TokenRefreshAPI = "https://northeurope.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	WestEuropeToken    TokenRefreshAPI = "https://westeurope.api.cognitive.microsoft.com/sts/v1.0/issueToken"
)

// serviceNameMappingString is a text placeholder that is present in the standard voices service name mapping. This is
// used to build the `name` attribute in the ssml/xml payload (see `voiceXML()`).
const serviceNameMappingString = "Microsoft Server Speech Text to Speech Voice"

// voiceXML renders the XML payload for the TTS api.
// For API reference see https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#sample-request-1
func voiceXML(speechText, description string, locale Region, gender Gender) string {

	return fmt.Sprintf(`<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s %s'>%s</voice></speak>`,
		locale, locale, gender, serviceNameMappingString, description, speechText)
}

// Synthesize returns a bytestream of the rendered text-to-speech in the target audio format. `speechText` is the string of
// text in which a user wishes to Synthesize, `region` is the language/locale, `gender` is the desired output voice
// and `audioOutput` captures the audio format.
func (az *AzureCSTextToSpeech) Synthesize(speechText string, region Region, gender Gender, audioOutput AudioOutput) ([]byte, error) {

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

	client := &http.Client{Timeout: synthesizeActionTimeout}
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

// RefreshToken fetches an updated token from the Azure cognitive speech/text services, or an error if unable to retrive.
// Each token is valid for a maximum of 10 minutes. Details for auth tokens are referenced at
// https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#authentication .
// Note: This does not need to be called by a client, since this automatically runs via a background go-routine (`startRefresher`)
func (az *AzureCSTextToSpeech) RefreshToken() error {
	request, _ := http.NewRequest(http.MethodPost, az.tokenAPI, nil)
	request.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	client := &http.Client{Timeout: tokenRefreshTimeout}

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

// startRefresher updates the authentication token on at a 9 minute interval. A channel is returned
// if the caller wishes to cancel the channel.
func (az *AzureCSTextToSpeech) startRefresher() chan bool {
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

// AzureCSTextToSpeech stores configuration and state information for the TTS client.
type AzureCSTextToSpeech struct {
	accessToken        string    // Token received from `TokenRefreshAPI`
	APIBase            string    // target `TextToSpeechAPI`
	SubscriptionKey    string    // API key for Azure's Congnitive Speech services
	TokenRefreshDoneCh chan bool // channel to stop the token refresh goroutine.
	tokenAPI           string    // Local endpoints for TokenRefreshAPI
}

// New returns an `AzureCSTextToSpeech` object and starts a background token refresh timer
func New(subscriptionKey string, api TextToSpeechAPI, token TokenRefreshAPI) (*AzureCSTextToSpeech, error) {

	az := &AzureCSTextToSpeech{
		APIBase:         string(api) + v1Path,
		SubscriptionKey: subscriptionKey,
		tokenAPI:        string(token)}

	// api requires that the token is refreshed every 10 mintutes.
	// We will do this task in the background every ~9 minutes.
	if err := az.RefreshToken(); err != nil {
		return nil, fmt.Errorf("failed to fetch initial token, err='%+v'", err)
	}
	az.TokenRefreshDoneCh = az.startRefresher()
	return az, nil
}
