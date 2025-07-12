package models

import (
	"time"

	"gorm.io/gorm"
)

type Website struct {
	gorm.Model
	URL         string `gorm:"uniqueIndex"`
	Name        string
	LastChecked time.Time
}

type Manga struct {
	gorm.Model
	Title         string
	Description   string
	WebsiteID     uint
	Website       Website
	Tags          []Tag `gorm:"many2many:manga_tags;"`
	Chapters      []Chapter
	LastChapter   string
	UpdateTime    time.Time
	EstimatedNext time.Time
	ExternalURL   string
	Author        string
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Chapter struct {
	gorm.Model
	MangaID     uint
	Number      uint
	Title       string
	ReleaseDate time.Time
	URL         string
}

type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	Email     string
	Password  string  // Hashed
	Favorites []Manga `gorm:"many2many:user_favorites;"`
	Bookmarks []Bookmark
	// Role string // For admin, add later if needed
}

type Bookmark struct {
	gorm.Model
	UserID  uint
	MangaID uint
	Chapter uint
}

type Notification struct {
	gorm.Model
	UserID  uint
	MangaID uint
	Message string
	SentAt  time.Time
	Read    bool // Add to track read status
}
