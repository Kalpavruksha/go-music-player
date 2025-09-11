# 🎵 Advanced Music Player Features - Implementation Examples

This document provides practical code examples for implementing the advanced features in your Go music player.

## 🎵 Playlists Implementation

### Backend - Playlist Model and API

Here's a complete example of how to implement playlists in your Go music player:

```go
// models/playlist.go
package models

import (
    "time"
)

type Playlist struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    IsPublic    bool      `json:"is_public"`
    UserID      uint      `json:"user_id"`
    Songs       []Song    `json:"songs" gorm:"many2many:playlist_songs;"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type PlaylistSong struct {
    PlaylistID uint `gorm:"primaryKey"`
    SongID     uint `gorm:"primaryKey"`
    Position   int  `gorm:"default:0"`
}
```

```go
// handlers/playlist_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
    "your-music-player/models"
)

func (h *Handler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
    var playlist models.Playlist
    if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := h.DB.Create(&playlist).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(playlist)
}

func (h *Handler) GetPlaylists(w http.ResponseWriter, r *http.Request) {
    var playlists []models.Playlist
    h.DB.Preload("Songs").Find(&playlists)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(playlists)
}

func (h *Handler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])
    
    var playlist models.Playlist
    if err := h.DB.Preload("Songs").First(&playlist, id).Error; err != nil {
        http.Error(w, "Playlist not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(playlist)
}

func (h *Handler) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    playlistID, _ := strconv.Atoi(vars["id"])
    
    var requestData map[string]int
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Invalid request data", http.StatusBadRequest)
        return
    }
    
    songID := requestData["song_id"]
    
    playlistSong := models.PlaylistSong{
        PlaylistID: uint(playlistID),
        SongID:     uint(songID),
    }
    
    if err := h.DB.Create(&playlistSong).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Song added to playlist"})
}
```

### Frontend - Playlist UI

```html
<!-- static/playlists.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Music Player - Playlists</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-900 text-white">
    <div class="container mx-auto p-6">
        <h1 class="text-3xl font-bold mb-6">My Playlists</h1>
        
        <!-- Create Playlist Form -->
        <div class="mb-8 p-4 bg-gray-800 rounded-lg">
            <h2 class="text-xl font-semibold mb-4">Create New Playlist</h2>
            <form id="create-playlist-form" class="space-y-4">
                <input type="text" id="playlist-name" placeholder="Playlist Name" 
                       class="w-full p-2 bg-gray-700 rounded">
                <textarea id="playlist-desc" placeholder="Description" 
                          class="w-full p-2 bg-gray-700 rounded"></textarea>
                <button type="submit" class="px-4 py-2 bg-purple-600 rounded hover:bg-purple-700">
                    Create Playlist
                </button>
            </form>
        </div>
        
        <!-- Playlists List -->
        <div id="playlists-container" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <!-- Playlists will be loaded here -->
        </div>
    </div>

    <script>
        // Load playlists on page load
        document.addEventListener('DOMContentLoaded', loadPlaylists);
        
        // Handle form submission
        document.getElementById('create-playlist-form').addEventListener('submit', createPlaylist);
        
        async function loadPlaylists() {
            try {
                const response = await fetch('/api/playlists');
                const playlists = await response.json();
                
                const container = document.getElementById('playlists-container');
                container.innerHTML = '';
                
                playlists.forEach(playlist => {
                    const playlistElement = document.createElement('div');
                    playlistElement.className = 'bg-gray-800 rounded-lg p-4 hover:bg-gray-700 transition';
                    playlistElement.innerHTML = `
                        <h3 class="text-xl font-semibold">${playlist.name}</h3>
                        <p class="text-gray-400 text-sm mb-2">${playlist.description || 'No description'}</p>
                        <p class="text-gray-500 text-xs">${playlist.songs.length} songs</p>
                        <button onclick="viewPlaylist(${playlist.id})" 
                                class="mt-3 px-3 py-1 bg-purple-600 rounded text-sm hover:bg-purple-700">
                            View Playlist
                        </button>
                    `;
                    container.appendChild(playlistElement);
                });
            } catch (error) {
                console.error('Error loading playlists:', error);
            }
        }
        
        async function createPlaylist(e) {
            e.preventDefault();
            
            const name = document.getElementById('playlist-name').value;
            const description = document.getElementById('playlist-desc').value;
            
            if (!name) {
                alert('Please enter a playlist name');
                return;
            }
            
            try {
                const response = await fetch('/api/playlists', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ name, description, is_public: false }),
                });
                
                if (response.ok) {
                    document.getElementById('playlist-name').value = '';
                    document.getElementById('playlist-desc').value = '';
                    loadPlaylists(); // Refresh the list
                }
            } catch (error) {
                console.error('Error creating playlist:', error);
            }
        }
        
        function viewPlaylist(id) {
            // Redirect to playlist detail page
            window.location.href = `/playlist.html?id=${id}`;
        }
    </script>
</body>
</html>
```

## 🔊 Equalizer Implementation

### Backend - Audio Processing

```go
// audio/equalizer.go
package audio

import (
    "github.com/u2takey/ffmpeg-go"
)

type EqualizerSettings struct {
    Frequency float64 `json:"frequency"`
    Width     float64 `json:"width"`
    Gain      float64 `json:"gain"`
}

func ApplyEqualizer(inputFile, outputFile string, settings EqualizerSettings) error {
    return ffmpeg.Input(inputFile).
        Filter("equalizer", ffmpeg.Args{
            "frequency": settings.Frequency,
            "width":     settings.Width,
            "gain":      settings.Gain,
        }).
        Output(outputFile).
        OverWriteOutput().
        Run()
}

// Apply multiple bands
func ApplyMultiBandEqualizer(inputFile, outputFile string, bands []EqualizerSettings) error {
    chain := ffmpeg.Input(inputFile)
    
    for _, band := range bands {
        chain = chain.Filter("equalizer", ffmpeg.Args{
            "frequency": band.Frequency,
            "width":     band.Width,
            "gain":      band.Gain,
        })
    }
    
    return chain.Output(outputFile).
        OverWriteOutput().
        Run()
}
```

### Frontend - Equalizer UI

```html
<!-- static/equalizer.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Music Player - Equalizer</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-900 text-white">
    <div class="container mx-auto p-6">
        <h1 class="text-3xl font-bold mb-6">Audio Equalizer</h1>
        
        <div class="bg-gray-800 rounded-lg p-6 mb-6">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-semibold">Equalizer Settings</h2>
                <button id="reset-eq" class="px-3 py-1 bg-gray-700 rounded hover:bg-gray-600">
                    Reset
                </button>
            </div>
            
            <!-- Equalizer Bands -->
            <div class="space-y-6">
                <div>
                    <div class="flex justify-between mb-2">
                        <span>Bass (60Hz)</span>
                        <span id="bass-value">0 dB</span>
                    </div>
                    <input type="range" id="bass" min="-12" max="12" value="0" 
                           class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer">
                </div>
                
                <div>
                    <div class="flex justify-between mb-2">
                        <span>Low Mids (250Hz)</span>
                        <span id="low-mids-value">0 dB</span>
                    </div>
                    <input type="range" id="low-mids" min="-12" max="12" value="0" 
                           class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer">
                </div>
                
                <div>
                    <div class="flex justify-between mb-2">
                        <span>High Mids (1kHz)</span>
                        <span id="high-mids-value">0 dB</span>
                    </div>
                    <input type="range" id="high-mids" min="-12" max="12" value="0" 
                           class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer">
                </div>
                
                <div>
                    <div class="flex justify-between mb-2">
                        <span>Presence (4kHz)</span>
                        <span id="presence-value">0 dB</span>
                    </div>
                    <input type="range" id="presence" min="-12" max="12" value="0" 
                           class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer">
                </div>
                
                <div>
                    <div class="flex justify-between mb-2">
                        <span>Brilliance (16kHz)</span>
                        <span id="brilliance-value">0 dB</span>
                    </div>
                    <input type="range" id="brilliance" min="-12" max="12" value="0" 
                           class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer">
                </div>
            </div>
            
            <div class="mt-6">
                <button id="apply-eq" class="w-full py-3 bg-purple-600 rounded hover:bg-purple-700">
                    Apply Equalizer
                </button>
            </div>
        </div>
        
        <!-- Visualizer -->
        <div class="bg-gray-800 rounded-lg p-6">
            <h2 class="text-xl font-semibold mb-4">Audio Visualizer</h2>
            <canvas id="visualizer" class="w-full h-32 bg-black rounded"></canvas>
        </div>
    </div>

    <script>
        // Get all slider elements
        const sliders = {
            bass: document.getElementById('bass'),
            'low-mids': document.getElementById('low-mids'),
            'high-mids': document.getElementById('high-mids'),
            presence: document.getElementById('presence'),
            brilliance: document.getElementById('brilliance')
        };
        
        // Get value display elements
        const valueDisplays = {
            bass: document.getElementById('bass-value'),
            'low-mids': document.getElementById('low-mids-value'),
            'high-mids': document.getElementById('high-mids-value'),
            presence: document.getElementById('presence-value'),
            brilliance: document.getElementById('brilliance-value')
        };
        
        // Update value displays when sliders change
        Object.keys(sliders).forEach(band => {
            sliders[band].addEventListener('input', function() {
                valueDisplays[band].textContent = this.value + ' dB';
            });
        });
        
        // Reset button
        document.getElementById('reset-eq').addEventListener('click', function() {
            Object.keys(sliders).forEach(band => {
                sliders[band].value = 0;
                valueDisplays[band].textContent = '0 dB';
            });
        });
        
        // Apply equalizer
        document.getElementById('apply-eq').addEventListener('click', async function() {
            const eqSettings = {
                bands: [
                    { frequency: 60, width: 1, gain: parseFloat(sliders.bass.value) },
                    { frequency: 250, width: 1, gain: parseFloat(sliders['low-mids'].value) },
                    { frequency: 1000, width: 1, gain: parseFloat(sliders['high-mids'].value) },
                    { frequency: 4000, width: 1, gain: parseFloat(sliders.presence.value) },
                    { frequency: 16000, width: 1, gain: parseFloat(sliders.brilliance.value) }
                ]
            };
            
            try {
                const response = await fetch('/api/audio/equalizer', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(eqSettings),
                });
                
                if (response.ok) {
                    alert('Equalizer settings applied!');
                } else {
                    alert('Failed to apply equalizer settings');
                }
            } catch (error) {
                console.error('Error applying equalizer:', error);
                alert('Error applying equalizer settings');
            }
        });
        
        // Audio visualizer
        const canvas = document.getElementById('visualizer');
        const ctx = canvas.getContext('2d');
        
        // Set canvas dimensions
        function resizeCanvas() {
            canvas.width = canvas.offsetWidth;
            canvas.height = canvas.offsetHeight;
        }
        
        window.addEventListener('resize', resizeCanvas);
        resizeCanvas();
        
        // Draw visualizer (simplified example)
        function drawVisualizer() {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            
            // Draw frequency bars
            const barCount = 64;
            const barWidth = canvas.width / barCount;
            
            for (let i = 0; i < barCount; i++) {
                const height = Math.random() * canvas.height;
                const hue = (i / barCount) * 360;
                
                ctx.fillStyle = `hsl(${hue}, 100%, 50%)`;
                ctx.fillRect(i * barWidth, canvas.height - height, barWidth - 2, height);
            }
            
            requestAnimationFrame(drawVisualizer);
        }
        
        drawVisualizer();
    </script>
</body>
</html>
```

## 🧠 Smart Recommendations

### Backend - Recommendation Engine

```go
// recommendations/engine.go
package recommendations

import (
    "math"
)

type UserHistory struct {
    UserID   uint `gorm:"index"`
    SongID   uint `gorm:"index"`
    PlayCount int
    LastPlayed int64
}

// Calculate similarity between two users based on listening history
func calculateUserSimilarity(user1History, user2History []UserHistory) float64 {
    // Create maps for easier lookup
    user1Map := make(map[uint]int)
    user2Map := make(map[uint]int)
    
    for _, h := range user1History {
        user1Map[h.SongID] = h.PlayCount
    }
    
    for _, h := range user2History {
        user2Map[h.SongID] = h.PlayCount
    }
    
    // Calculate cosine similarity
    var dotProduct, magnitude1, magnitude2 float64
    
    allSongs := make(map[uint]bool)
    for songID := range user1Map {
        allSongs[songID] = true
    }
    for songID := range user2Map {
        allSongs[songID] = true
    }
    
    for songID := range allSongs {
        playCount1 := float64(user1Map[songID])
        playCount2 := float64(user2Map[songID])
        
        dotProduct += playCount1 * playCount2
        magnitude1 += playCount1 * playCount1
        magnitude2 += playCount2 * playCount2
    }
    
    if magnitude1 == 0 || magnitude2 == 0 {
        return 0
    }
    
    return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

// Get recommendations for a user
func GetUserRecommendations(userID uint, db *gorm.DB) ([]Song, error) {
    var userHistory []UserHistory
    db.Where("user_id = ?", userID).Find(&userHistory)
    
    // Find similar users
    var allUsers []uint
    db.Model(&UserHistory{}).Distinct("user_id").Pluck("user_id", &allUsers)
    
    type SimilarUser struct {
        UserID     uint
        Similarity float64
    }
    
    var similarUsers []SimilarUser
    
    for _, otherUserID := range allUsers {
        if otherUserID == userID {
            continue
        }
        
        var otherUserHistory []UserHistory
        db.Where("user_id = ?", otherUserID).Find(&otherUserHistory)
        
        similarity := calculateUserSimilarity(userHistory, otherUserHistory)
        if similarity > 0.1 { // Threshold for similarity
            similarUsers = append(similarUsers, SimilarUser{otherUserID, similarity})
        }
    }
    
    // Sort by similarity
    sort.Slice(similarUsers, func(i, j int) bool {
        return similarUsers[i].Similarity > similarUsers[j].Similarity
    })
    
    // Get top recommendations
    recommendationMap := make(map[uint]float64) // songID -> weighted score
    
    for _, similarUser := range similarUsers[:min(10, len(similarUsers))] {
        var similarUserHistory []UserHistory
        db.Where("user_id = ?", similarUser.UserID).Find(&similarUserHistory)
        
        for _, history := range similarUserHistory {
            // Check if user already listened to this song
            alreadyListened := false
            for _, userHist := range userHistory {
                if userHist.SongID == history.SongID {
                    alreadyListened = true
                    break
                }
            }
            
            if !alreadyListened {
                // Weight by similarity and play count
                score := similarUser.Similarity * float64(history.PlayCount)
                recommendationMap[history.SongID] += score
            }
        }
    }
    
    // Convert to slice and sort by score
    type Recommendation struct {
        SongID uint
        Score  float64
    }
    
    var recommendations []Recommendation
    for songID, score := range recommendationMap {
        recommendations = append(recommendations, Recommendation{songID, score})
    }
    
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].Score > recommendations[j].Score
    })
    
    // Get top 20 recommendations
    var songIDs []uint
    for i, rec := range recommendations {
        if i >= 20 {
            break
        }
        songIDs = append(songIDs, rec.SongID)
    }
    
    // Fetch song details
    var songs []Song
    if len(songIDs) > 0 {
        db.Where("id IN ?", songIDs).Find(&songs)
    }
    
    return songs, nil
}
```

## ⏯ Gapless Playback + Crossfade

### Frontend Implementation

```javascript
// static/js/player-enhancements.js
class AdvancedPlayer {
    constructor() {
        this.audio = new Audio();
        this.nextAudio = new Audio();
        this.isPlaying = false;
        this.currentSong = null;
        this.nextSong = null;
        this.crossfadeDuration = 5000; // 5 seconds
        this.gaplessEnabled = true;
        this.crossfadeEnabled = true;
        
        this.setupEventListeners();
    }
    
    setupEventListeners() {
        // Setup audio event listeners
        this.audio.addEventListener('timeupdate', () => this.onTimeUpdate());
        this.audio.addEventListener('ended', () => this.onSongEnded());
        this.audio.addEventListener('loadstart', () => this.onLoadStart());
        this.audio.addEventListener('canplay', () => this.onCanPlay());
    }
    
    // Gapless playback implementation
    onTimeUpdate() {
        if (!this.gaplessEnabled || !this.nextSong) return;
        
        // Preload next song when we're close to the end
        if (this.audio.duration - this.audio.currentTime < 10) {
            if (this.nextAudio.src !== this.nextSong.url) {
                this.nextAudio.src = this.nextSong.url;
                this.nextAudio.load();
            }
        }
        
        // Switch to next song with no gap
        if (this.audio.duration - this.audio.currentTime < 0.05) {
            this.switchToNextSong();
        }
    }
    
    switchToNextSong() {
        if (!this.nextSong) return;
        
        // Swap audio elements
        const tempAudio = this.audio;
        this.audio = this.nextAudio;
        this.nextAudio = tempAudio;
        
        // Update current song
        this.currentSong = this.nextSong;
        this.nextSong = null;
        
        // Continue playback
        this.audio.play();
        
        // Dispatch event for UI update
        this.dispatchEvent('songchange', this.currentSong);
    }
    
    // Crossfade implementation
    async playWithCrossfade(song) {
        if (this.isPlaying && this.crossfadeEnabled) {
            // Fade out current song
            await this.fadeOut(this.audio, this.crossfadeDuration);
        }
        
        // Set up new song
        this.audio.src = song.url;
        await this.audio.play();
        
        // Fade in new song
        if (this.crossfadeEnabled) {
            await this.fadeIn(this.audio, this.crossfadeDuration);
        }
        
        this.currentSong = song;
        this.isPlaying = true;
    }
    
    fadeOut(audio, duration) {
        return new Promise((resolve) => {
            const startVolume = audio.volume;
            const startTime = Date.now();
            
            const fadeOutInterval = setInterval(() => {
                const elapsed = Date.now() - startTime;
                const progress = Math.min(elapsed / duration, 1);
                
                audio.volume = startVolume * (1 - progress);
                
                if (progress >= 1) {
                    clearInterval(fadeOutInterval);
                    resolve();
                }
            }, 50);
        });
    }
    
    fadeIn(audio, duration) {
        return new Promise((resolve) => {
            const targetVolume = 1;
            const startTime = Date.now();
            
            // Start at 0 volume
            audio.volume = 0;
            
            const fadeInInterval = setInterval(() => {
                const elapsed = Date.now() - startTime;
                const progress = Math.min(elapsed / duration, 1);
                
                audio.volume = targetVolume * progress;
                
                if (progress >= 1) {
                    clearInterval(fadeInInterval);
                    resolve();
                }
            }, 50);
        });
    }
    
    // Set next song for gapless playback
    setNextSong(song) {
        this.nextSong = song;
        
        // Preload next song
        if (this.gaplessEnabled) {
            this.nextAudio.src = song.url;
            this.nextAudio.load();
        }
    }
    
    // Event dispatcher
    dispatchEvent(type, detail) {
        const event = new CustomEvent(type, { detail });
        document.dispatchEvent(event);
    }
}

// Initialize advanced player
const player = new AdvancedPlayer();

// Example usage
document.addEventListener('songchange', (e) => {
    console.log('Now playing:', e.detail);
    // Update UI with new song info
});
```

## 🖼 Album Art + Lyrics Sync

### Backend - Album Art Fetching

```go
// artwork/fetcher.go
package artwork

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

type LastFmResponse struct {
    Results struct {
        TrackMatches struct {
            Track []struct {
                Name   string `json:"name"`
                Artist string `json:"artist"`
                Image  []struct {
                    Size string `json:"size"`
                    URL  string `json:"#text"`
                } `json:"image"`
            } `json:"track"`
        } `json:"trackmatches"`
    } `json:"results"`
}

func FetchAlbumArtFromLastFM(artist, track, apiKey string) (string, error) {
    baseURL := "http://ws.audioscrobbler.com/2.0/"
    
    params := url.Values{}
    params.Add("method", "track.search")
    params.Add("track", track)
    params.Add("artist", artist)
    params.Add("api_key", apiKey)
    params.Add("format", "json")
    
    fullURL := baseURL + "?" + params.Encode()
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(fullURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    var response LastFmResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return "", err
    }
    
    if len(response.Results.TrackMatches.Track) > 0 {
        track := response.Results.TrackMatches.Track[0]
        if len(track.Image) > 0 {
            // Get the largest image (usually the last one)
            return track.Image[len(track.Image)-1].URL, nil
        }
    }
    
    return "", fmt.Errorf("no artwork found")
}
```

### Frontend - Lyrics Synchronization

```html
<!-- static/lyrics.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Music Player - Lyrics</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-900 text-white">
    <div class="container mx-auto p-6">
        <h1 class="text-3xl font-bold mb-6">Now Playing with Lyrics</h1>
        
        <div class="flex flex-col lg:flex-row gap-6">
            <!-- Album Art -->
            <div class="lg:w-1/3">
                <div class="bg-gray-800 rounded-lg p-6">
                    <div class="bg-gray-700 rounded-xl w-full aspect-square flex items-center justify-center mb-4">
                        <img id="album-art" src="" alt="Album Art" class="rounded-xl w-full h-full object-cover hidden">
                        <div id="album-art-placeholder" class="text-6xl text-gray-600">
                            <i class="fas fa-compact-disc"></i>
                        </div>
                    </div>
                    <h2 id="song-title" class="text-2xl font-bold text-center">Song Title</h2>
                    <p id="song-artist" class="text-gray-400 text-center">Artist Name</p>
                </div>
            </div>
            
            <!-- Lyrics -->
            <div class="lg:w-2/3">
                <div class="bg-gray-800 rounded-lg p-6">
                    <h2 class="text-xl font-semibold mb-4">Lyrics</h2>
                    <div id="lyrics-container" class="h-96 overflow-y-auto space-y-2">
                        <!-- Lyrics will be loaded here -->
                    </div>
                </div>
                
                <!-- Audio Controls -->
                <div class="bg-gray-800 rounded-lg p-6 mt-6">
                    <div class="flex items-center justify-center space-x-4">
                        <button id="prev-btn" class="text-2xl text-gray-400 hover:text-white">
                            <i class="fas fa-step-backward"></i>
                        </button>
                        <button id="play-btn" class="text-3xl text-white bg-purple-600 rounded-full w-16 h-16 flex items-center justify-center hover:bg-purple-700">
                            <i class="fas fa-play"></i>
                        </button>
                        <button id="next-btn" class="text-2xl text-gray-400 hover:text-white">
                            <i class="fas fa-step-forward"></i>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        class LyricsSync {
            constructor() {
                this.lyrics = [];
                this.currentLineIndex = -1;
                this.audio = new Audio();
                this.isPlaying = false;
                
                this.setupEventListeners();
            }
            
            setupEventListeners() {
                this.audio.addEventListener('timeupdate', () => this.syncLyrics());
                this.audio.addEventListener('play', () => this.onPlay());
                this.audio.addEventListener('pause', () => this.onPause());
                
                document.getElementById('play-btn').addEventListener('click', () => this.togglePlay());
                document.getElementById('prev-btn').addEventListener('click', () => this.previous());
                document.getElementById('next-btn').addEventListener('click', () => this.next());
            }
            
            loadSong(song) {
                // Set audio source
                this.audio.src = song.audioUrl;
                
                // Update UI
                document.getElementById('song-title').textContent = song.title;
                document.getElementById('song-artist').textContent = song.artist;
                
                // Load album art
                if (song.albumArt) {
                    document.getElementById('album-art').src = song.albumArt;
                    document.getElementById('album-art').classList.remove('hidden');
                    document.getElementById('album-art-placeholder').classList.add('hidden');
                }
                
                // Load lyrics
                this.loadLyrics(song.lyricsUrl);
            }
            
            async loadLyrics(lyricsUrl) {
                try {
                    const response = await fetch(lyricsUrl);
                    const lyricsText = await response.text();
                    
                    // Parse LRC format
                    this.lyrics = this.parseLrc(lyricsText);
                    this.renderLyrics();
                } catch (error) {
                    console.error('Error loading lyrics:', error);
                }
            }
            
            parseLrc(lrcText) {
                const lines = lrcText.split('\n');
                const lyrics = [];
                
                const timeRegex = /\[(\d{2}):(\d{2})\.(\d{2})\](.*)/;
                
                for (const line of lines) {
                    const match = line.match(timeRegex);
                    if (match) {
                        const minutes = parseInt(match[1]);
                        const seconds = parseInt(match[2]);
                        const centiseconds = parseInt(match[3]);
                        const text = match[4];
                        
                        const time = minutes * 60 + seconds + centiseconds / 100;
                        
                        lyrics.push({
                            time: time,
                            text: text
                        });
                    }
                }
                
                return lyrics;
            }
            
            renderLyrics() {
                const container = document.getElementById('lyrics-container');
                container.innerHTML = '';
                
                this.lyrics.forEach((line, index) => {
                    const lineElement = document.createElement('div');
                    lineElement.className = 'p-2 rounded text-center text-gray-400 transition-colors';
                    lineElement.textContent = line.text;
                    lineElement.dataset.index = index;
                    container.appendChild(lineElement);
                });
            }
            
            syncLyrics() {
                if (this.lyrics.length === 0) return;
                
                const currentTime = this.audio.currentTime;
                let newLineIndex = -1;
                
                // Find current line
                for (let i = 0; i < this.lyrics.length; i++) {
                    if (this.lyrics[i].time <= currentTime) {
                        newLineIndex = i;
                    } else {
                        break;
                    }
                }
                
                // Update highlighting if line changed
                if (newLineIndex !== this.currentLineIndex) {
                    this.highlightLine(newLineIndex);
                    this.currentLineIndex = newLineIndex;
                }
            }
            
            highlightLine(index) {
                // Remove previous highlighting
                const previousLine = document.querySelector('.bg-purple-600');
                if (previousLine) {
                    previousLine.classList.remove('bg-purple-600', 'text-white');
                    previousLine.classList.add('text-gray-400');
                }
                
                // Highlight new line
                if (index >= 0 && index < this.lyrics.length) {
                    const lineElement = document.querySelector(`[data-index="${index}"]`);
                    if (lineElement) {
                        lineElement.classList.remove('text-gray-400');
                        lineElement.classList.add('bg-purple-600', 'text-white');
                        
                        // Scroll to line
                        lineElement.scrollIntoView({
                            behavior: 'smooth',
                            block: 'center'
                        });
                    }
                }
            }
            
            togglePlay() {
                if (this.isPlaying) {
                    this.audio.pause();
                } else {
                    this.audio.play();
                }
            }
            
            onPlay() {
                this.isPlaying = true;
                document.getElementById('play-btn').innerHTML = '<i class="fas fa-pause"></i>';
            }
            
            onPause() {
                this.isPlaying = false;
                document.getElementById('play-btn').innerHTML = '<i class="fas fa-play"></i>';
            }
            
            previous() {
                // Implement previous song logic
            }
            
            next() {
                // Implement next song logic
            }
        }
        
        // Initialize lyrics sync
        const lyricsSync = new LyricsSync();
        
        // Example usage
        // lyricsSync.loadSong({
        //     title: "Song Title",
        //     artist: "Artist Name",
        //     albumArt: "/api/artwork/123",
        //     audioUrl: "/api/song/123",
        //     lyricsUrl: "/api/lyrics/123"
        // });
    </script>
</body>
</html>
```

This implementation provides practical examples for all the advanced features you requested. Each feature is broken down into backend and frontend components with working code examples that can be integrated into your existing Go music player.