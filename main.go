package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"DiscordBot/voice"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/Andreychik32/ytdl"
	
)

var prefix = "!"

var voiceConnections[] Voice
var stopChannel chan bool
var queue [] Song
var l1 = "https://www.youtube.com/watch?v=H6z1QjXqTCw"
var l2 = "https://www.youtube.com/watch?v=z2JkCXAZZnc"
var description = `You must make a decision that you are going to move on.It wont happen automatically.You will have to rise up and say,'I don't care how hard this is,I don't care how disappointed l am,I'm not going to let this get the best of me.I'm moving on with my life."-Joel Osteen

`
func main() {
	// Creating a new bot with the specified token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	stopChannel = make(chan bool)

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

	//voiceChannel := findVoiceChannelID(g, m)

	var commandArgs []string = strings.Split(m.Content, " ")
	if commandArgs[0] == prefix + "play" {
		voiceChannel := findVoiceChannelID(g, m)
		voiceConnections = append(voiceConnections,connectToVoiceChannel(s, m.GuildID, voiceChannel))
	} else if commandArgs[0] == prefix + "youtube" {
		voiceChannel := findVoiceChannelID(g, m)
		voiceConnections = append(voiceConnections,connectToVoiceChannel(s, m.GuildID, voiceChannel))
		go playYoutubeLink(s, commandArgs[1], m.GuildID, voiceChannel, m.ChannelID)
	} else if commandArgs[0] == prefix + "disconnect" {
		go disconnectFromVoiceChannel(m.GuildID)
	} else if commandArgs[0] == prefix + "skip" {
		voiceChannel := findVoiceChannelID(g, m)
		stopChannel <- true
		playYoutubeLink(s, l1, m.GuildID, voiceChannel, m.ChannelID)
	} 
}

func addSong(song Song) {
	fmt.Println("Added to the Queue")
	queue = append(queue, song)
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

func checkForDoubleVoiceConnection(guild string, channel string) {
	for index, voice := range voiceConnections {
		if voice.Guild == guild {
			voiceConnections = append(voiceConnections[:index], voiceConnections[index+1:]...)
		}
	}
}

func playYoutubeLink(bot *discordgo.Session, link string, guild string, channel string, textChannel string){
	ctx := context.Background()
	videoInfo, err := ytdl.GetVideoInfo(ctx,link)
	if err != nil {
		fmt.Println(err)
		return // Returning to avoid crash when video informations could not be found
	}
	var messageEmbed discordgo.MessageEmbed
	//p := &messageEmbed.Title
	//fmt.Println(p)
	messageEmbed.Title = videoInfo.Title
	messageEmbed.Description = description
	p := &messageEmbed
	//fmt.Println(sendMessage)
	fmt.Printf("%+v", messageEmbed)

	bot.ChannelMessageSendEmbed(textChannel, p)
	for _, format := range videoInfo.Formats {
		if format.AudioEncoding == "opus" || format.AudioEncoding == "aac" || format.AudioEncoding == "vorbis" {
			url := format.URL
			go playAudioFile(url, guild, channel, "youtube")
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

func nextSong() {
	if len(queue) > 0 {
		go playAudioFile(queue[0].Link, queue[0].Guild, queue[0].Channel, queue[0].Type)
		queue = append(queue[:0], queue[1:]...)
	} 
}

func playAudioFile(file string, guild string, channel string, linkType string) {
	voiceConnection, index := findVoiceConnection(guild, channel)
	switch voiceConnection.PlayerStatus {
	case false:
		voiceConnections[index].PlayerStatus = true
		dgvoice.PlayAudioFile(voiceConnection.VoiceConnection, file, stopChannel)
		voiceConnections[index].PlayerStatus = false
	case true:
		addSong(Song{
			Link:    file,
			Type:    linkType,
			Guild:   guild,
			Channel: channel,
		})
	}
}

func disconnectFromVoiceChannel(guild string) {
	for index, voice := range voiceConnections {
		if voice.Guild == guild {
			_ = voice.VoiceConnection.Disconnect()
			stopChannel <- true
			voiceConnections = append(voiceConnections[:index], voiceConnections[index+1:]...)
		}
	}
}