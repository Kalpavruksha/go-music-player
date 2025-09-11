// sync_client.js - Frontend JavaScript for multi-device sync
// This would be included in your HTML file

class MusicSyncClient {
    constructor() {
        this.ws = null;
        this.userID = this.generateUserID();
        this.isConnected = false;
        this.currentSong = null;
        this.currentPosition = 0;
        this.isPlaying = false;
        
        this.initializeWebSocket();
        this.setupAudioEventListeners();
    }
    
    // Generate a unique user ID
    generateUserID() {
        // In a real app, this would come from authentication
        return 'user_' + Math.random().toString(36).substr(2, 9);
    }
    
    // Initialize WebSocket connection
    initializeWebSocket() {
        const wsURL = `ws://${window.location.host}/ws/${this.userID}`;
        
        this.ws = new WebSocket(wsURL);
        
        this.ws.onopen = () => {
            console.log('Connected to sync server');
            this.isConnected = true;
            this.sendPresence();
        };
        
        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleSyncMessage(message);
            } catch (error) {
                console.error('Error parsing sync message:', error);
            }
        };
        
        this.ws.onclose = () => {
            console.log('Disconnected from sync server');
            this.isConnected = false;
            
            // Attempt to reconnect
            setTimeout(() => {
                this.initializeWebSocket();
            }, 5000);
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }
    
    // Handle incoming sync messages
    handleSyncMessage(message) {
        // Ignore messages from ourselves
        if (message.user_id === this.userID) {
            return;
        }
        
        switch (message.type) {
            case 'play':
                this.syncPlay(message.song_id, message.position, message.timestamp);
                break;
            case 'pause':
                this.syncPause(message.position, message.timestamp);
                break;
            case 'seek':
                this.syncSeek(message.position, message.timestamp);
                break;
            case 'volume':
                this.syncVolume(message.volume);
                break;
            case 'playlist':
                this.syncPlaylist(message.playlist_id);
                break;
        }
    }
    
    // Sync play action
    syncPlay(songID, position, timestamp) {
        // Only sync if the message is recent (within 5 seconds)
        const messageAge = (Date.now() / 1000) - timestamp;
        if (messageAge > 5) {
            return;
        }
        
        console.log(`Sync play: ${songID} at ${position}s`);
        
        // Update UI to show what friend is playing
        this.showFriendPlaying(songID, message.user_id);
        
        // In a real app, you might want to:
        // - Show notification
        // - Update friend activity feed
        // - Allow joining the session
    }
    
    // Sync pause action
    syncPause(position, timestamp) {
        const messageAge = (Date.now() / 1000) - timestamp;
        if (messageAge > 5) {
            return;
        }
        
        console.log(`Sync pause at ${position}s`);
    }
    
    // Send play event
    sendPlay(songID, position) {
        if (!this.isConnected) return;
        
        const message = {
            type: 'play',
            user_id: this.userID,
            song_id: songID,
            position: position,
            timestamp: Math.floor(Date.now() / 1000)
        };
        
        this.ws.send(JSON.stringify(message));
    }
    
    // Send pause event
    sendPause(position) {
        if (!this.isConnected) return;
        
        const message = {
            type: 'pause',
            user_id: this.userID,
            position: position,
            timestamp: Math.floor(Date.now() / 1000)
        };
        
        this.ws.send(JSON.stringify(message));
    }
    
    // Send seek event
    sendSeek(position) {
        if (!this.isConnected) return;
        
        const message = {
            type: 'seek',
            user_id: this.userID,
            position: position,
            timestamp: Math.floor(Date.now() / 1000)
        };
        
        this.ws.send(JSON.stringify(message));
    }
    
    // Send volume change
    sendVolume(volume) {
        if (!this.isConnected) return;
        
        const message = {
            type: 'volume',
            user_id: this.userID,
            volume: volume,
            timestamp: Math.floor(Date.now() / 1000)
        };
        
        this.ws.send(JSON.stringify(message));
    }
    
    // Send playlist change
    sendPlaylist(playlistID) {
        if (!this.isConnected) return;
        
        const message = {
            type: 'playlist',
            user_id: this.userID,
            playlist_id: playlistID,
            timestamp: Math.floor(Date.now() / 1000)
        };
        
        this.ws.send(JSON.stringify(message));
    }
    
    // Setup audio event listeners
    setupAudioEventListeners() {
        // These would be connected to your actual audio player
        document.addEventListener('songPlay', (e) => {
            this.currentSong = e.detail.songID;
            this.currentPosition = e.detail.position;
            this.isPlaying = true;
            this.sendPlay(this.currentSong, this.currentPosition);
        });
        
        document.addEventListener('songPause', (e) => {
            this.currentPosition = e.detail.position;
            this.isPlaying = false;
            this.sendPause(this.currentPosition);
        });
        
        document.addEventListener('songSeek', (e) => {
            this.currentPosition = e.detail.position;
            this.sendSeek(this.currentPosition);
        });
        
        document.addEventListener('volumeChange', (e) => {
            this.sendVolume(e.detail.volume);
        });
    }
    
    // Show friend playing notification
    showFriendPlaying(songID, friendID) {
        // Create or update notification
        const notificationID = `friend-${friendID}`;
        let notification = document.getElementById(notificationID);
        
        if (!notification) {
            notification = document.createElement('div');
            notification.id = notificationID;
            notification.className = 'fixed bottom-4 right-4 bg-purple-600 text-white p-4 rounded-lg shadow-lg z-50';
            document.body.appendChild(notification);
        }
        
        // Get song info (in real app, you'd fetch this from API)
        const songInfo = this.getSongInfo(songID);
        
        notification.innerHTML = `
            <div class="flex items-center">
                <div class="w-10 h-10 rounded-full bg-purple-800 flex items-center justify-center mr-3">
                    <i class="fas fa-user"></i>
                </div>
                <div>
                    <div class="font-semibold">Friend Playing</div>
                    <div class="text-sm">${songInfo.title}</div>
                </div>
                <button onclick="this.closeFriendNotification('${friendID}')" 
                        class="ml-4 text-purple-200 hover:text-white">
                    <i class="fas fa-times"></i>
                </button>
            </div>
        `;
        
        // Auto-hide after 5 seconds
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 5000);
    }
    
    // Get song information (placeholder)
    getSongInfo(songID) {
        // In a real app, this would fetch from your API
        return {
            title: `Song ${songID}`,
            artist: 'Unknown Artist',
            album: 'Unknown Album'
        };
    }
    
    // Send presence heartbeat
    sendPresence() {
        if (!this.isConnected) return;
        
        const message = {
            type: 'presence',
            user_id: this.userID,
            timestamp: Math.floor(Date.now() / 1000),
            playing: this.isPlaying,
            song_id: this.currentSong
        };
        
        this.ws.send(JSON.stringify(message));
        
        // Send presence every 30 seconds
        setTimeout(() => {
            this.sendPresence();
        }, 30000);
    }
}

// Close friend notification
function closeFriendNotification(friendID) {
    const notification = document.getElementById(`friend-${friendID}`);
    if (notification) {
        notification.parentNode.removeChild(notification);
    }
}

// Initialize the sync client when the page loads
document.addEventListener('DOMContentLoaded', () => {
    window.musicSync = new MusicSyncClient();
});

// Example usage in your player controls:
/*
// When user plays a song
document.getElementById('play-button').addEventListener('click', () => {
    const songID = getCurrentSongID();
    const position = getCurrentPosition();
    window.musicSync.sendPlay(songID, position);
});

// When user pauses
document.getElementById('pause-button').addEventListener('click', () => {
    const position = getCurrentPosition();
    window.musicSync.sendPause(position);
});

// When user seeks
document.getElementById('progress-bar').addEventListener('input', (e) => {
    const position = e.target.value;
    window.musicSync.sendSeek(position);
});
*/