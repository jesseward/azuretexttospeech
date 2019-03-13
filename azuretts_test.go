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

	assert.NoError(t, err)
	assert.Equal(t, az.SubscriptionKey, az.accessToken)
}
