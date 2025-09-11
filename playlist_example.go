// playlist_example.go - Example implementation of playlist functionality
// This is a conceptual example showing how to extend the current music player

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Song represents a music file
type Song struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	File      string    `json:"file"`
	CreatedAt time.Time `json:"created_at"`
}

// Playlist represents a collection of songs
type Playlist struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	Songs       []Song    `gorm:"many2many:playlist_songs;" json:"songs"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlaylistSong is the join table for many-to-many relationship
type PlaylistSong struct {
	PlaylistID uint `gorm:"primaryKey"`
	SongID     uint `gorm:"primaryKey"`
	Position   int  `gorm:"default:0"`
}

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Initialize database
	db, err := gorm.Open(sqlite.Open("music.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Song{}, &Playlist{}, &PlaylistSong{})

	// API Routes
	app.Get("/api/songs", getSongs(db))
	app.Post("/api/songs", createSong(db))

	app.Get("/api/playlists", getPlaylists(db))
	app.Post("/api/playlists", createPlaylist(db))
	app.Get("/api/playlists/:id", getPlaylist(db))
	app.Put("/api/playlists/:id", updatePlaylist(db))
	app.Delete("/api/playlists/:id", deletePlaylist(db))

	app.Post("/api/playlists/:id/songs", addSongToPlaylist(db))
	app.Delete("/api/playlists/:id/songs/:songId", removeSongFromPlaylist(db))

	// Start server
	log.Fatal(app.Listen(":8080"))
}

// getSongs returns all songs
func getSongs(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var songs []Song
		db.Find(&songs)
		return c.JSON(songs)
	}
}

// createSong adds a new song
func createSong(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		song := new(Song)
		if err := c.BodyParser(song); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		db.Create(&song)
		return c.JSON(song)
	}
}

// getPlaylists returns all playlists
func getPlaylists(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var playlists []Playlist
		db.Preload("Songs").Find(&playlists)
		return c.JSON(playlists)
	}
}

// createPlaylist creates a new playlist
func createPlaylist(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		playlist := new(Playlist)
		if err := c.BodyParser(playlist); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		db.Create(&playlist)
		return c.JSON(playlist)
	}
}

// getPlaylist returns a specific playlist with its songs
func getPlaylist(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var playlist Playlist
		result := db.Preload("Songs").First(&playlist, id)
		if result.Error != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
		}
		return c.JSON(playlist)
	}
}

// updatePlaylist updates an existing playlist
func updatePlaylist(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var playlist Playlist
		result := db.First(&playlist, id)
		if result.Error != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
		}

		if err := c.BodyParser(&playlist); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		db.Save(&playlist)
		return c.JSON(playlist)
	}
}

// deletePlaylist removes a playlist
func deletePlaylist(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result := db.Delete(&Playlist{}, id)
		if result.Error != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
		}
		return c.SendString("Playlist deleted")
	}
}

// addSongToPlaylist adds a song to a playlist
func addSongToPlaylist(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		playlistID := c.Params("id")

		// Parse request body for song ID
		var requestData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		songID, ok := requestData["song_id"].(float64)
		if !ok {
			return c.Status(400).JSON(fiber.Map{"error": "Missing or invalid song_id"})
		}

		// Check if playlist exists
		var playlist Playlist
		if err := db.First(&playlist, playlistID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
		}

		// Check if song exists
		var song Song
		if err := db.First(&song, int(songID)).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Song not found"})
		}

		// Add song to playlist
		playlistSong := PlaylistSong{
			PlaylistID: uint(playlistID),
			SongID:     uint(songID),
		}

		if err := db.Create(&playlistSong).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to add song to playlist"})
		}

		return c.JSON(fiber.Map{"message": "Song added to playlist"})
	}
}

// removeSongFromPlaylist removes a song from a playlist
func removeSongFromPlaylist(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		playlistID := c.Params("id")
		songID := c.Params("songId")

		// Remove song from playlist
		result := db.Where("playlist_id = ? AND song_id = ?", playlistID, songID).
			Delete(&PlaylistSong{})

		if result.Error != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to remove song from playlist"})
		}

		if result.RowsAffected == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Song not found in playlist"})
		}

		return c.JSON(fiber.Map{"message": "Song removed from playlist"})
	}
}
