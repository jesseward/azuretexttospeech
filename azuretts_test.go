package azuretexttospeech

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoiceXML(t *testing.T) {
	expect := "<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Female' name='Microsoft Server Speech Text to Speech Voice (en-US, ZiraRUS)'>Microsoft Speech Service Text-to-Speech API</voice></speak>"
	gender := Female
	region := EnUS
	description := localeToGender[localeGender{region, gender}]
	assert.Equal(t, expect, voiceXML("Microsoft Speech Service Text-to-Speech API", description, region, gender))
}

func TestSynthesize(t *testing.T) {
	az := &AzureCSTextToSpeech{SubscriptionKey: "SYS64738", accessToken: "SYS49152"}

	// payload should be nil and err should be true, since DeCH + Female is not a valid combination
	payload, err := az.Synthesize("test-speech", DeCH, Female, RIFF8Bit8kHzMonoPCM)
	assert.Error(t, err, "should raise an error")
	assert.Nil(t, payload, "payload should be nil")

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("SYS4096"))
		}),
	)
	defer ts.Close()
	az.APIBase = ts.URL
	payload, err = az.Synthesize("SYS4096", EnUS, Female, RIFF8Bit8kHzMonoPCM)
	assert.NoError(t, err)
	assert.Equal(t, payload, []byte("SYS4096"))
}

// TestRefreshToken validates logic for fetching of the RefreshToken
func TestRefreshToken(t *testing.T) {
	az := &AzureCSTextToSpeech{SubscriptionKey: "ThisIsMySubscriptionKeyAndToBeToken"}

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// return the SubscriptionKey as the token, for test case only.
			w.Write([]byte(az.SubscriptionKey))
		}),
	)
	defer ts.Close()

	az.tokenAPI = ts.URL
	err := az.RefreshToken()

	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, az.SubscriptionKey, az.accessToken, "values should be equal")
}
