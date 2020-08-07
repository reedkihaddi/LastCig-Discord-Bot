package main

import "github.com/bwmarrin/discordgo"

// Voice discord.
type Voice struct {
	VoiceConnection *discordgo.VoiceConnection
	Channel         string
	Guild           string
	PlayerStatus    bool
}

// Song discord.
type Song struct {
	Link    string
	Type    string
	Guild   string
	Channel string
}

