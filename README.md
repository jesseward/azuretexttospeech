AzureTextToSpeech Client
---
This package provides a client for Azure's Cognitive Services (speech services) Text To Speech API. Users of the client
can specify the lanaguage (`Region` type), a string containing the desired text to speak as well as the gender (`Gender` type) in which the audiofile should be rendered. The library fetches the audio rendered in the format of your choice (see `AudioOutput` types for supported formats).

API documents of interest
* Text to speech [Azure pricing details](https://azure.microsoft.com/en-us/pricing/details/cognitive-services/speech-services/). Note there is a *free* tier available.
* Text to speech, speech services [API specifications](https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-apis#text-to-speech-api).

Requirements
---
A Cognitive Services (kind=Speech Services) API key is required to access the URL. This service can be enabled at the Azure portal.

Howto
---
The following will digitize the string `64 BASIC BYTES FREE. READY.`, using the en-US locale, rending with a female voice. The output file format is a 16khz 32kbit single channel MP3 audio file.

```
# See AzureCognitiveServicesAPI and AzureCognitiveServicesToken types for list of endpoints and regions.
azureSpeech, _ := tts.New("YOUR-API-KEY", WestUS2, WestUS2Token)
payload, _err_ := azureSpeech.Digitize(
    "64 BASIC BYTES FREE. READY.",
    EnUS, // Region type
    tts.Female, // Gender type
    tts.Audio16khz32kbitrateMonoMp3) // AudioOutput type
```
