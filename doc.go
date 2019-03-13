/*
package azuretexttospeech provides a client for Azure's Cognitive Services (speech services) Text To Speech API. Users of the client
can specify the locale (lanaguage), text in which to speak/digitize as well as the gender in which the gender should be rendered.

For Azure pricing see https://azure.microsoft.com/en-us/pricing/details/cognitive-services/speech-services/ . Note
there is a limited use *FREE* tier available.

Documentation for the TTS Cognitive Services API can be found at https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#text-to-speech-api

An API key is required to access the Azure API.

USAGE

	import (
		"io/ioutil"
		"os"

		tts "github.com/jesseward/azuretexttospeech"
	)

	func main() {
		// create a key for "Cognitive Services" (kind=SpeechServices). Once the key is available
		// in the azure portal, push it into an environment variable.
		// By default the free tier keys are served out of West US2
		az, err := tts.New(os.Getenv("AZUREKEY"), tts.WestUS2, tts.WestUS2Token)
		if err != nil {
			panic(err)
		}
		defer close(az.TokenRefreshDoneCh)

		// Digitize a text string using the enUS locale, female voice and specify the
		// audio format of a 16Khz, 32kbit mp3 file.
		b, err := az.Synthesize(
			"64 BASIC BYTES FREE. READY.",
			tts.EnUS,
			tts.Female,
			tts.Audio16khz32kbitrateMonoMp3)

		if err != nil {
			panic(err)
		}

		// write the results to disk.
		err = ioutil.WriteFile("audio.mp3", b, 0644)
		if err != nil {
			panic(err)
		}
	}

*/
package azuretexttospeech
