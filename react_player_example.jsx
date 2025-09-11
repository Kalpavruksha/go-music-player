// This is a conceptual example of how the React player component might look
// using modern libraries like Howler.js and Wavesurfer.js

import React, { useState, useEffect, useRef } from 'react';
import { Howl, Howler } from 'howler';
import WaveSurfer from 'wavesurfer.js';

const MusicPlayer = () => {
  const [songs, setSongs] = useState([]);
  const [currentSong, setCurrentSong] = useState(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [volume, setVolume] = useState(0.8);
  const [progress, setProgress] = useState(0);
  
  const waveformRef = useRef(null);
  const wavesurfer = useRef(null);
  const sound = useRef(null);

  // Fetch songs from API
  useEffect(() => {
    const fetchSongs = async () => {
      try {
        const response = await fetch('/api/songs');
        const data = await response.json();
        setSongs(data);
      } catch (error) {
        console.error('Error fetching songs:', error);
      }
    };

    fetchSongs();
  }, []);

  // Initialize waveform
  useEffect(() => {
    if (waveformRef.current) {
      wavesurfer.current = WaveSurfer.create({
        container: waveformRef.current,
        waveColor: '#8B5CF6',
        progressColor: '#EC4899',
        height: 80,
        responsive: true,
        normalize: true,
      });

      wavesurfer.current.on('audioprocess', () => {
        if (wavesurfer.current) {
          setProgress(wavesurfer.current.getCurrentTime());
        }
      });

      wavesurfer.current.on('seek', () => {
        if (wavesurfer.current) {
          setProgress(wavesurfer.current.getCurrentTime());
        }
      });
    }

    return () => {
      if (wavesurfer.current) {
        wavesurfer.current.destroy();
      }
    };
  }, []);

  // Play song
  const playSong = (song) => {
    // Stop current song if playing
    if (sound.current) {
      sound.current.unload();
    }

    // Load new song with Howler
    sound.current = new Howl({
      src: [`/api/song/${song.id}`],
      volume: volume,
      html5: true,
      onplay: () => {
        setIsPlaying(true);
        // Load waveform
        if (wavesurfer.current) {
          wavesurfer.current.load(`/api/song/${song.id}`);
        }
      },
      onpause: () => {
        setIsPlaying(false);
      },
      onend: () => {
        setIsPlaying(false);
      },
      onload: () => {
        console.log('Song loaded');
      }
    });

    setCurrentSong(song);
    sound.current.play();
  };

  // Toggle play/pause
  const togglePlay = () => {
    if (sound.current) {
      if (isPlaying) {
        sound.current.pause();
      } else {
        sound.current.play();
      }
    }
  };

  // Set volume
  const handleVolumeChange = (e) => {
    const newVolume = parseFloat(e.target.value);
    setVolume(newVolume);
    if (sound.current) {
      sound.current.volume(newVolume);
    }
  };

  // Seek to position
  const handleSeek = (e) => {
    if (wavesurfer.current && sound.current) {
      const clickX = e.nativeEvent.offsetX;
      const width = e.target.clientWidth;
      const duration = sound.current.duration();
      const seekTime = (clickX / width) * duration;
      
      sound.current.seek(seekTime);
      wavesurfer.current.seekTo(clickX / width);
      setProgress(seekTime);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-black text-white">
      {/* Header */}
      <header className="p-6 border-b border-gray-800">
        <h1 className="text-3xl font-bold bg-gradient-to-r from-purple-500 to-pink-500 bg-clip-text text-transparent">
          Advanced Go Music Player
        </h1>
      </header>

      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Song Library */}
          <div className="lg:col-span-2">
            <div className="bg-gray-800/50 rounded-2xl p-6 backdrop-blur-sm">
              <h2 className="text-2xl font-bold mb-6 flex items-center">
                <i className="fas fa-music mr-3 text-purple-500"></i>
                Song Library
              </h2>
              
              <div className="space-y-3 max-h-[500px] overflow-y-auto">
                {songs.map((song) => (
                  <div 
                    key={song.id}
                    className={`flex items-center p-4 rounded-xl cursor-pointer transition-all hover:bg-gray-700/50 ${
                      currentSong?.id === song.id ? 'bg-purple-900/30 border border-purple-500/30' : 'bg-gray-900/50'
                    }`}
                    onClick={() => playSong(song)}
                  >
                    <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center mr-4">
                      <i className="fas fa-music text-white"></i>
                    </div>
                    <div className="flex-1">
                      <h3 className="font-medium">{song.title}</h3>
                      <p className="text-gray-400 text-sm">{song.artist}</p>
                    </div>
                    <button className="text-gray-400 hover:text-white">
                      <i className="fas fa-play"></i>
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Player */}
          <div className="lg:col-span-1">
            <div className="bg-gray-800/50 rounded-2xl p-6 backdrop-blur-sm sticky top-6">
              <h2 className="text-2xl font-bold mb-6 text-center">
                Now Playing
              </h2>
              
              {/* Album Art */}
              <div className="bg-gradient-to-br from-purple-900/30 to-pink-900/30 rounded-xl p-8 mb-6 flex items-center justify-center">
                <div className="bg-gray-700 rounded-full w-48 h-48 flex items-center justify-center">
                  <i className="fas fa-compact-disc text-6xl text-gray-600"></i>
                </div>
              </div>

              {/* Song Info */}
              <div className="text-center mb-6">
                <h3 className="text-xl font-bold truncate">
                  {currentSong ? currentSong.title : 'No song selected'}
                </h3>
                <p className="text-gray-400">
                  {currentSong ? currentSong.artist : 'Select a song to play'}
                </p>
              </div>

              {/* Waveform */}
              <div 
                ref={waveformRef}
                className="mb-6 bg-gray-900/50 rounded-lg h-20 cursor-pointer"
                onClick={handleSeek}
              ></div>

              {/* Progress */}
              <div className="flex justify-between text-sm text-gray-400 mb-6">
                <span>{new Date(progress * 1000).toISOString().substr(14, 5)}</span>
                <span>{currentSong ? '3:45' : '0:00'}</span>
              </div>

              {/* Controls */}
              <div className="flex justify-center items-center space-x-6 mb-6">
                <button className="text-gray-400 hover:text-white">
                  <i className="fas fa-random"></i>
                </button>
                <button className="text-gray-400 hover:text-white">
                  <i className="fas fa-step-backward"></i>
                </button>
                <button 
                  className="bg-gradient-to-r from-purple-600 to-pink-600 rounded-full w-14 h-14 flex items-center justify-center hover:from-purple-500 hover:to-pink-500"
                  onClick={togglePlay}
                >
                  <i className={`fas ${isPlaying ? 'fa-pause' : 'fa-play'}`}></i>
                </button>
                <button className="text-gray-400 hover:text-white">
                  <i className="fas fa-step-forward"></i>
                </button>
                <button className="text-gray-400 hover:text-white">
                  <i className="fas fa-redo"></i>
                </button>
              </div>

              {/* Volume */}
              <div className="flex items-center space-x-3">
                <i className="fas fa-volume-down text-gray-400"></i>
                <input 
                  type="range" 
                  min="0" 
                  max="1" 
                  step="0.01" 
                  value={volume}
                  onChange={handleVolumeChange}
                  className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-purple-500"
                />
                <i className="fas fa-volume-up text-gray-400"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MusicPlayer;