package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

//const token = "NzQwODE5MjI5MzU5NDcyNjUw.XyujrA.SPXhJLd1YvvQYZCEigwG7YCJf0A"
const token = "NzQwOTk5MzIxMTQ2NjIyMDEz.XyxLZQ.EMbSESn4Syv4UFz92HNi5_NGc08"
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
	
	// channel, err := s.State.Channel(m.ChannelID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// guild, err := s.State.Guild(channel.GuildID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(guild)

	// guild, err := s.State.Guild(m.GuildID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// for _, guild := range s.State.Guilds{
	// 	channels, _ := s.GuildChannels(guild.ID)
	// 	for _, c := range channels {
	// 		if c.Type == discordgo.ChannelTypeGuildVoice{
	// 		fmt.Printf("Name of channel is: %s \n", c.Name)
	// 		}
	// 	}
	// }

	// if m.Author.ID == s.State.User.ID{
	// 	return
	// }

	// if m.Content == "ping" {
	// 	guild, err := s.State.Guild(m.GuildID)
	// 	if err != nil {
	// 		fmt.Println("Error")
	// 	}
	// 	fmt.Println(guild)
	//findVoiceChannel(s, m.GuildID, m.ChannelID, m.Author.ID)
	// 	s.ChannelMessageSend(m.ChannelID, "pong")

	// }

	// if m.Content == "pong" {
	// 	s.ChannelMessageSend(m.ChannelID, "ping")
	// }

	// if m.Author.ID == s.State.User.ID {
	// 	return
	// }
	// fmt.Printf("Channel ID is: %s ",m.ChannelID)
	// fmt.Printf("Guild ID is: %s\n",m.GuildID)
	// //fmt.Println(s.State.Channel(m.ChannelID))
	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Guild not found.")
		// Could not find guild.
		return
	}
	voiceChannel := findVoiceChannelID(g, m)
	fmt.Println("CHANNEL ID:::::: ",voiceChannel)
	fmt.Printf("%+v\n",g)
	// for _, vs := range g.VoiceStates {
	// 	if vs.UserID == m.Author.ID {
	// 		fmt.Println(1)
	// 	}
	// }

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

func findVoiceChannel(bot *discordgo.Session, guild string, channel string, author string) {
	//var voiceConnection Voice
	//var index int
	//index = 1
	// fmt.Println(bot.State.Guild)
	
	// channels, _ := bot.GuildChannels(guild)
	// voice,err := bot.State.Guild(guild)
	// if err != nil {
	// 	fmt.Println("Error")
	// }
	// fmt.Println(voice.Name)
	// for _, v := range voice.VoiceStates{
	// 	fmt.Println(v.UserID)
	// }

	for _, g := range bot.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == author {
				fmt.Println(v.ChannelID)
			}
		}
	}

	// for _, c := range channels {
	// 	if c.Type == discordgo.ChannelTypeGuildVoice{
	// 		fmt.Printf("ID of channel is: %s \n ", c.ID)
	// 	}
	// }
	
	//return voiceConnection, index
}

