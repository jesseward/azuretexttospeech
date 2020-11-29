package azuretexttospeech

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const textToSpeechAPI = "https://%s.tts.speech.microsoft.com/cognitiveservices/v1"
const tokenRefreshAPI = "https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken"

// synthesizeActionTimeout is the amount of time the http client will wait for a response during Synthesize request
const synthesizeActionTimeout = time.Second * 30

// tokenRefreshTimeout is the amount of time the http client will wait during the token refresh action.
const tokenRefreshTimeout = time.Second * 15

// TTSApiXMLPayload templates the payload required for API.
// See: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#sample-request
const TTSApiXMLPayload = "<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'>%s</voice></speak>"

// Region references the locations of the availability of standard voices.
// See https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/regions#standard-voices
type Region int

const (
	// Azure regions and their endpoints that support the Text To Speech service.
	RegionAustraliaEast Region = iota
	RegionBrazilSouth
	RegionCanadaCentral
	RegionCentralUS
	RegionEastAsia
	RegionEastUS
	RegionEastUS2
	RegionFranceCentral
	RegionIndiaCentral
	RegionJapanEast
	RegionJapanWest
	RegionKoreaCentral
	RegionNorthCentralUS
	RegionNorthEurope
	RegionSouthCentralUS
	RegionSoutheastAsia
	RegionUKSouth
	RegionWestEurope
	RegionWestUS
	RegionWestUS2
)

func (t Region) String() string {
	return [...]string{
		"australiaeast",
		"brazilsouth",
		"canadacentral",
		"centralus",
		"eastasia",
		"eastus",
		"eastus2",
		"francecentral",
		"indiacentral",
		"japaneast",
		"japanwest",
		"koreacentral",
		"northcentralus",
		"northeurope",
		"southcentralus",
		"southeastasia",
		"uksouth",
		"westeurope",
		"westus",
		"westus2",
	}[t]

}

// SynthesizeWithContext returns a bytestream of the rendered text-to-speech in the target audio format. `speechText` is the string of
// text in which a user wishes to Synthesize, `region` is the language/locale, `gender` is the desired output voice
// and `audioOutput` captures the audio format.
func (az *AzureCSTextToSpeech) SynthesizeWithContext(ctx context.Context, speechText string, locale Locale, gender Gender, audioOutput AudioOutput) ([]byte, error) {

	description, ok := az.RegionVoiceMap[supportedVoices{gender, locale}]
	if !ok {
		return nil, fmt.Errorf("unable to to locate RegionVoiceMap{region=%s, gender=%s} pair", locale, gender)
	}

	v := voiceXML(speechText, description, locale, gender)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, az.textToSpeechURL, bytes.NewBufferString(v))
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Microsoft-OutputFormat", fmt.Sprint(audioOutput))
	request.Header.Set("Content-Type", "application/ssml+xml")
	request.Header.Set("Authorization", "Bearer "+az.accessToken)
	request.Header.Set("User-Agent", "azuretts")

	client := &http.Client{}
	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// list of acceptable response status codes
	// see: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#http-status-codes-1
	switch response.StatusCode {
	case http.StatusOK:
		// The request was successful; the response body is an audio file.
		return ioutil.ReadAll(response.Body)
	case http.StatusBadRequest:
		return nil, fmt.Errorf("%d - A required parameter is missing, empty, or null. Or, the value passed to either a required or optional parameter is invalid. A common issue is a header that is too long", response.StatusCode)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%d - The request is not authorized. Check to make sure your subscription key or token is valid and in the correct region", response.StatusCode)
	case http.StatusRequestEntityTooLarge:
		return nil, fmt.Errorf("%d - The SSML input is longer than 1024 characters", response.StatusCode)
	case http.StatusUnsupportedMediaType:
		return nil, fmt.Errorf("%d - It's possible that the wrong Content-Type was provided. Content-Type should be set to application/ssml+xml", response.StatusCode)
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("%d - You have exceeded the quota or rate of requests allowed for your subscription", response.StatusCode)
	case http.StatusBadGateway:
		return nil, fmt.Errorf("%d - Network or server-side issue. May also indicate invalid headers", response.StatusCode)
	}

	return nil, fmt.Errorf("%d - received unexpected HTTP status code", response.StatusCode)
}

// Synthesize directs to SynthesizeWithContext
func (az *AzureCSTextToSpeech) Synthesize(speechText string, locale Locale, gender Gender, audioOutput AudioOutput) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), synthesizeActionTimeout)
	defer cancel()
	return az.SynthesizeWithContext(ctx, speechText, locale, gender, audioOutput)
}

// voiceXML renders the XML payload for the TTS api.
// For API reference see https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#sample-request
func voiceXML(speechText, description string, locale Locale, gender Gender) string {
	return fmt.Sprintf(TTSApiXMLPayload, locale, locale, gender, description, speechText)
}

// refreshToken fetches an updated token from the Azure cognitive speech/text services, or an error if unable to retrive.
// Each token is valid for a maximum of 10 minutes. Details for auth tokens are referenced at
// https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#authentication .
// Note: This does not need to be called by a client, since this automatically runs via a background go-routine (`startRefresher`)
func (az *AzureCSTextToSpeech) refreshToken() error {
	request, _ := http.NewRequest(http.MethodPost, az.tokenRefreshURL, nil)
	request.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	client := &http.Client{Timeout: tokenRefreshTimeout}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code; received http status=%s", response.Status)
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
				err := az.refreshToken()
				if err != nil {
					log.Printf("failed to refresh token, %v", err)
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
	accessToken         string // is the auth token received from `TokenRefreshAPI`. Used in the Authorization: Bearer header.
	RegionVoiceMap      RegionVoiceMap
	SubscriptionKey     string    // API key for Azure's Congnitive Speech services
	TokenRefreshDoneCh  chan bool // channel to stop the token refresh goroutine.
	tokenRefreshURL     string
	voiceServiceListURL string
	textToSpeechURL     string
}

// New returns an `AzureCSTextToSpeech` object and starts a background token refresh timer
func New(subscriptionKey string, region Region) (*AzureCSTextToSpeech, error) {
	az := &AzureCSTextToSpeech{
		SubscriptionKey: subscriptionKey,
	}

	az.textToSpeechURL = fmt.Sprintf(textToSpeechAPI, region)
	az.tokenRefreshURL = fmt.Sprintf(tokenRefreshAPI, region)
	az.voiceServiceListURL = fmt.Sprintf(voiceListAPI, region)

	// api requires that the token is refreshed every 10 mintutes.
	// We will do this task in the background every ~9 minutes.
	if err := az.refreshToken(); err != nil {
		return nil, fmt.Errorf("failed to fetch initial token, %v", err)
	}

	m, err := az.buildVoiceToRegionMap()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch voice-map, %v", err)
	}
	az.RegionVoiceMap = m

	az.TokenRefreshDoneCh = az.startRefresher()
	return az, nil
}
