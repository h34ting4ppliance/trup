package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

const removeUsage = "remove <amount> <@user>"

func remove(ctx *Context, args []string) {
	if len(args) < 3 {
		ctx.Reply("not enough arguments.")
		return
	}

	number, err := strconv.Atoi(args[1])
	if err != nil {
		ctx.ReportError("the first argument must be a number", err)
		return
	} else if number > 100 || number < 2 {
		ctx.Reply("the first argument must be comprised between 2 and 100")
		return
	}

	from := parseMention(args[2])
	if from == "" {
		ctx.Reply("the second argument must be a user mention")
		return
	}

	var toDelete []string
	before := ctx.Message.ID
	tooOldThreshold := time.Hour * 24 * 14
Outer:
	for i := 1; i < 10; i++ {
		messages, err := ctx.Session.ChannelMessages(ctx.Message.ChannelID, 100, before, "", "")
		if err != nil {
			ctx.ReportError("error", err)
			continue
		}
		for _, message := range messages {
			created, _ := discordgo.SnowflakeTimestamp(message.ID)
			duration := time.Now().Sub(created)
			if duration > tooOldThreshold {
				break Outer
			}
			if message.Author.ID == from {
				toDelete = append(toDelete, message.ID)
			}
			if len(toDelete) >= number {
				break Outer
			}
		}
		before = messages[len(messages)-1].ID
		if len(messages) < 100 {
			break
		}
	}
	ctx.Session.ChannelMessagesBulkDelete(ctx.Message.ChannelID, toDelete)
	ctx.Reply(fmt.Sprintf("Deleted %d messages", len(toDelete)))
}