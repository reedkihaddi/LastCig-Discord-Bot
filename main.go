package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/Andreychik32/ytdl"
	
)

const token = "NzQwODE5MjI5MzU5NDcyNjUw.XyujrA.SPXhJLd1YvvQYZCEigwG7YCJf0A"
//const token = "NzQwOTk5MzIxMTQ2NjIyMDEz.XyxLZQ.EMbSESn4Syv4UFz92HNi5_NGc08"
var prefix = "!"

var voiceConnections[] Voice

func main() {
	// Creating a new bot with the specified token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	go dg.AddHandler(createMessage)


	err = dg.Open()
	if err != nil {
		fmt.Println("Error connecting to the bot.")
		return
	}
	fmt.Println("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func createMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Guild not found.")
		return
	}

	voiceChannel := findVoiceChannelID(g, m)

	var commandArgs []string = strings.Split(m.Content, " ")
	if commandArgs[0] == prefix + "play" {
		voiceConnections = append(voiceConnections,connectToVoiceChannel(s, m.GuildID, voiceChannel))
	}
	if commandArgs[0] == prefix + "youtube" {
		go playYoutubeLink(commandArgs[1], m.GuildID, voiceChannel)
	}

	fmt.Printf("Voice Channel used is in: %s\n", voiceChannel)

}

func findVoiceChannelID(guild *discordgo.Guild, message *discordgo.MessageCreate) string {
	var channelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == message.Author.ID {
			channelID = vs.ChannelID
		}
	}
	return channelID
}

func connectToVoiceChannel(bot *discordgo.Session, guild string, channel string) Voice {
	vs, err := bot.ChannelVoiceJoin(guild, channel, false, true)
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

func playYoutubeLink(link string, guild string, channel string){
	ctx := context.Background()
	videoInfo, err := ytdl.GetVideoInfo(ctx,link)
	fmt.Printf("%+v\n",videoInfo)
	if err != nil {
		fmt.Println(err)
		return // Returning to avoid crash when video informations could not be found
	}
	client := ytdl.DefaultClient
	for _, format := range videoInfo.Formats {
		//fmt.Printf("%+v\n",format)
		if format.AudioEncoding == "opus" || format.AudioEncoding == "aac" || format.AudioEncoding == "vorbis" {
			url := format.URL
			fmt.Println(url)
		}	
	}
	file, err := os.Create("2" + ".mp4")
	if err != nil {
		fmt.Println("1",err)
	}
	defer file.Close()

	err = client.Download(ctx, videoInfo, videoInfo.Formats[0], file)
	if err != nil {
		fmt.Println("1",err)
	}
}