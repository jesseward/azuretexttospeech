package main

import (
	"io/ioutil"
	"os"

	tts "github.com/jesseward/azuretexttospeech"
)

func main() {
	// create a key for "Cognitive Services" (kind=SpeechServices). Once the key is available
	// in the Azure portal, push it into an environment variable (export AZUREKEY=SYS64738).
	// By default the free tier keys are served out of West US2
	var apiKey string
	if apiKey = os.Getenv("AZUREKEY"); apiKey == "" {
		panic("Please set your AZUREKEY environment variable")
	}
	az, err := tts.New(apiKey, tts.WestUS2, tts.WestUS2Token)
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
	// send results to disk.
	err = ioutil.WriteFile("audio.mp3", b, 0644)
	if err != nil {
		panic(err)
	}

}
