package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"DiscordBot/voice"
	"os"
	"os/signal"
	"encoding/json"
	"syscall"
	"strings"
	"github.com/Andreychik32/ytdl"
	"math/rand"
	"time"
	"strconv"
	"io/ioutil"
)


// Prefix for the command which bot searches for.
var prefix = "!lc"
// BOT token.
const token = "NzQwODE5MjI5MzU5NDcyNjUw.XyujrA.UDisrwBRMmg8x9f52UePP_IozV0"
// A list to track currently connected voice channels.
var voiceConnections[] Voice
// A struct that contains list of songs.
var songs SongList
// A channel for stopping songs.
var stopChannel chan bool


func main() {

	file, _ := ioutil.ReadFile("songs.json")
	_ = json.Unmarshal([]byte(file), &songs)

	// Creating a new bot with the specified token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Creating the channel.
	stopChannel = make(chan bool)

	// Add handler to scan messages.
	go dg.AddHandler(createMessage)

	// Connect to the bot.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error connecting to the bot.")
		return
	}

	fmt.Println("Bot is now running.")

	// To keep bot running.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}


func createMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Find the guild from where the message came from.
	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Guild not found.")
		return
	}

	// To check the message contents.
	var commandArgs []string = strings.Split(m.Content, " ")
	

	if commandArgs[0] == "!LastCig" {
		
		x1 := rand.NewSource(time.Now().UnixNano())
		y1 := rand.New(x1) 
		n := strconv.Itoa(y1.Intn(50))
		voiceChannel := findVoiceChannelID(g, m)
		voiceConnections = append(voiceConnections,connectToVoiceChannel(s, m.GuildID, voiceChannel))
		go playYoutubeLink(s, songs[n].Link, m.GuildID, voiceChannel, m.ChannelID, n)

	} else if commandArgs[0] == prefix && commandArgs[1] == "d" {
		
		go disconnectFromVoiceChannel(m.GuildID)

	} else if commandArgs[0] == prefix && commandArgs[1] == "skip" {
		// To get a random number. 
		x1 := rand.NewSource(time.Now().UnixNano())
		y1 := rand.New(x1) 
		n := strconv.Itoa(y1.Intn(50))
		voiceChannel := findVoiceChannelID(g, m)
		// Stop the current audio.
		stopChannel <- true
		// Play a new random song.
		go playYoutubeLink(s, songs[n].Link, m.GuildID, voiceChannel, m.ChannelID, n)
	
	} else if commandArgs[0] == "play" {

		voiceChannel := findVoiceChannelID(g, m)
		voiceConnections = append(voiceConnections,connectToVoiceChannel(s, m.GuildID, voiceChannel))
		go playYoutubeLink(s, commandArgs[1], m.GuildID, voiceChannel, m.ChannelID, "3")

	} else if commandArgs[0] == prefix && commandArgs[1] == "help" {

		var messageEmbed discordgo.MessageEmbed
		var messageEmbedFooter discordgo.MessageEmbedFooter
		messageEmbed.Title = "Last Cigarette Bot"
		messageEmbed.Description = "**A simple music bot made purely in Golang.**\n\n**!LastCig** to play an endless playlist of Last Cigarette songs.\n\n**!lc skip** to skip the current song.\n\n**!lc disconnect** to disconnect the bot from the voice channel.\n\n"
		messageEmbedFooter.Text = "https://github.com/reedkihaddi/LastCig-Discord-Bot"
		messageEmbed.Footer = &messageEmbedFooter
		p := &messageEmbed

		s.ChannelMessageSendEmbed(m.ChannelID, p)
	}
}


// Find the voice channel ID of the author.
func findVoiceChannelID(guild *discordgo.Guild, message *discordgo.MessageCreate) string {
	var channelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == message.Author.ID {
			channelID = vs.ChannelID
		}
	}
	return channelID
}


// Connect to the voice channel and return Voice struct.
func connectToVoiceChannel(bot *discordgo.Session, guild string, channel string) Voice {
	vs, err := bot.ChannelVoiceJoin(guild, channel, false, true)

	checkForDoubleVoiceConnection(guild, channel)

	if err != nil {
		fmt.Println(err)
	}
	return Voice{
		VoiceConnection: vs,
		Channel:         channel,
		Guild:           guild,
		PlayerStatus:    false,
	}
}


// Check if bot is already present in the voice channel.
func checkForDoubleVoiceConnection(guild string, channel string) {
	for index, voice := range voiceConnections {
		if voice.Guild == guild {
			voiceConnections = append(voiceConnections[:index], voiceConnections[index+1:]...)
		}
	}
}


// Play youtube link.
func playYoutubeLink(bot *discordgo.Session, link string, guild string, channel string, textChannel string, n string){
	ctx := context.Background()
	
	// Get video info from link. Title, description etc...
	videoInfo, err := ytdl.GetVideoInfo(ctx,link)
	durationVideo := videoInfo.Duration.String()
	if err != nil {
		fmt.Println(err)
		return // Returning to avoid crash when video informations could not be found
	}

	// Check available formats for the link.
	for _, format := range videoInfo.Formats {
		if format.AudioEncoding == "opus" || format.AudioEncoding == "aac" || format.AudioEncoding == "vorbis" {
			url := format.URL
			//fmt.Println(url)
			// Send the file to play on Discord.

			go playAudioFile(bot, url, guild, channel, "youtube", textChannel, n, durationVideo)

			return 
		}	
	}

}


func findVoiceConnection(guild string, channel string) (Voice, int) {
	var voiceConnection Voice
	var index int
	for i, vc := range voiceConnections {
		if vc.Guild == guild {
			voiceConnection = vc
			index = i
		}
	}
	return voiceConnection, index
}


// Plays the audio file in Discord.
func playAudioFile(bot *discordgo.Session, file string, guild string, channel string, linkType string, textChannel string, n string, length string) {

	// Find the Voice Connection.
	voiceConnection, index := findVoiceConnection(guild, channel)

	switch voiceConnection.PlayerStatus {
	case false:
		
		voiceConnections[index].PlayerStatus = true
		fmt.Println(voiceConnections)
		
		var messageEmbed discordgo.MessageEmbed
		var messageEmbedFooter discordgo.MessageEmbedFooter
		messageEmbed.Title = songs[n].Title
		messageEmbed.Description = songs[n].Describe
		messageEmbed.Color = 15158332
		messageEmbed.URL = songs[n].Link
		messageEmbedFooter.Text = "Duration: " + length
		messageEmbed.Footer = &messageEmbedFooter
		p := &messageEmbed

		bot.ChannelMessageSendEmbed(textChannel, p)

		dgvoice.PlayAudioFile(voiceConnection.VoiceConnection, file, stopChannel)

		// Generate a random number to play a new song.
		x1 := rand.NewSource(time.Now().UnixNano())
		y1 := rand.New(x1) 
		n := strconv.Itoa(y1.Intn(50))

		// Play a new song.
		go playYoutubeLink(bot, songs[n].Link, guild, channel, textChannel, n)
		voiceConnections[index].PlayerStatus = false
		fmt.Println(voiceConnections)
	}
}

// Disconnects the bot from the voice 
func disconnectFromVoiceChannel(guild string) {
	for index, voice := range voiceConnections {
		if voice.Guild == guild {
			_ = voice.VoiceConnection.Disconnect()
			stopChannel <- true
			voiceConnections = append(voiceConnections[:index], voiceConnections[index+1:]...)
		}
	}
}