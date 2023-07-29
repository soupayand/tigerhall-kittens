package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID       int64
	Username string
	Password string
	Email    string
}

type Claims struct {
	UserID   int64  `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email_id,omitempty"`
	jwt.StandardClaims
}

type Animal struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Variant     string `json:"variant,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Description string `json:"description,omitempty"`
}

type Image struct {
	FileName string `json:"filename,omitempty"`
	Type     string `json:"file_type,omitempty"`
	data     []byte `json:"image_data,omitempty"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Sighting struct {
	ID           int64 `json:"id,omitempty"`
	Animal       Animal
	Image        Image
	Reporter     User
	LastLocation Point
	LastSeen     time.Time `json:"last_seen"`
}
