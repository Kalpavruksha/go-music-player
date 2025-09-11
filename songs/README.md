# Music Files Directory

Place your MP3/WAV/OGG/FLAC files in this directory.

## Supported Formats
- MP3 (.mp3)
- WAV (.wav)
- OGG (.ogg)
- FLAC (.flac)

## Adding Music Files

1. **Copy your music files** to this directory
2. **Refresh the music player** in your browser
3. **Your songs will appear** in the library automatically

## Sample Files

For testing purposes, you can create small sample files:

### Using ffmpeg (if available):
```bash
# Create a 10-second silent MP3
ffmpeg -f lavfi -i anullsrc=r=44100:cl=mono -t 10 -q:a 9 -acodec libmp3lame sample.mp3

# Create a 5-second sine wave
ffmpeg -f lavfi -i "sine=frequency=1000:duration=5" -acodec libmp3lame sample2.mp3
```

### Download Sample Files
Sample music files are available in the [Releases](https://github.com/Kalpavruksha/go-music-player/releases) section of this repository.

## File Naming
- Avoid special characters in filenames
- Use alphanumeric characters, spaces, hyphens, and underscores
- Example: "Artist - Song Title.mp3"