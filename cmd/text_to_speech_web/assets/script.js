function ttsPlayer() {
    return {
        text: 'Today is a wonderful day to build something people love!',
        voice: 'alloy',
        isFetching: false,
        isPlaying: false,
        isLoaded: false,
        textChanged() {
            this.isLoaded = false;
        },
        voiceChanged() {
            this.isLoaded = false;
        },
        async togglePlayPause() {
            const audio = this.$refs.audio;

            if (!this.isLoaded) {
                try {
                    const payload = new FormData();
                    payload.append('text', this.text);
                    payload.append('voice', this.voice);

                    this.isFetching = true;
                    const response = await fetch('/tts', {
                        method: 'POST',
                        body: payload
                    });
                    this.isFetching = false;

                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }

                    const blob = await response.blob();
                    audio.src = URL.createObjectURL(blob);
                    this.isLoaded = true;
                } catch (error) {
                    console.error('Error fetching TTS audio:', error);
                    return;
                }
            }

            if (audio.paused) {
                audio.play();
                this.isPlaying = true;
            } else {
                audio.pause();
                this.isPlaying = false;
            }
        },
        audioEnded() {
            this.isPlaying = false;
        }
    }
}
