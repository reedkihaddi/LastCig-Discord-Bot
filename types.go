package main

import "github.com/bwmarrin/discordgo"

// Voice discord.
type Voice struct {
	VoiceConnection *discordgo.VoiceConnection
	Channel         string
	Guild           string
	PlayerStatus    bool
}