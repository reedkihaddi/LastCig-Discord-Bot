package main

import "github.com/bwmarrin/discordgo"

// Voice discord.
type Voice struct {
	VoiceConnection *discordgo.VoiceConnection
	Channel         string
	Guild           string
	PlayerStatus    bool
}

// SongList check.
type SongList map[string]SongData

// SongID discord.
type SongID struct {
	ID SongData
}

// SongData struct.
type SongData struct {
	Title string
	Link string
	Describe string
}