package main

import (
	"fmt"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
	"log"
)

func sendEmbed(ctx *dgc.Ctx, title string, message string) *discordgo.Message {
	// create embed to send
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.Event.Author.Username,
			IconURL: ctx.Event.Author.AvatarURL("1024"),
		},
		Title:       title,
		Color:       0x607CA3,
		Description: message,
	}

	// send the Embed
	dcMessage, err := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embed)
	if err != nil {
		log.Printf("Something went wrong sending message to discord: %s", err)
	}

	return dcMessage
}

// Send the message to discord within the `fwew` layout of an embed.
func sendDiscordMessage(ctx *dgc.Ctx, message string) {
	// create title from executed command
	title := ctx.Command.Name
	arguments := ctx.Arguments.Raw()
	if arguments != "" {
		title += " " + arguments
	}

	sendEmbed(ctx, title, message)
}

type message struct {
	message *discordgo.Message
	title   string
	curPage *int
	pages   []string
}

var messages = map[string]message{}

func sendDiscordMessagePaginated(ctx *dgc.Ctx, pages []string) {
	// create title from executed command with pages count
	titleSimple := ctx.Command.Name
	arguments := ctx.Arguments.Raw()
	if arguments != "" {
		titleSimple += " " + arguments
	}

	var title string
	// add pages to
	if len(pages) > 1 {
		title = titleSimple + fmt.Sprintf(" (Page %d/%d)", 1, len(pages))
	}

	// post first page
	dcMessage := sendEmbed(ctx, title, pages[0])
	session := ctx.Session

	if len(pages) > 1 {
		// add arrows as reaction to pagination
		session.MessageReactionAdd(dcMessage.ChannelID, dcMessage.ID, "⬅️")
		session.MessageReactionAdd(dcMessage.ChannelID, dcMessage.ID, "➡️")

		// save message so pagination can work
		p := 1
		messages[dcMessage.ChannelID+":"+dcMessage.ID] = message{
			message: dcMessage,
			title:   titleSimple,
			pages:   pages,
			curPage: &p,
		}
	}
}