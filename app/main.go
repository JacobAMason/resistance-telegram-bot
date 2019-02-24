// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	// "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func messageContainsCommandMatcher(command string) func(Update) bool {
	return func(update Update) bool {
		return messageContainsCommand(strings.ToLower(update.Message.Text), strings.ToLower(command))
	}
}

func messageContainsCommand(message string, command string) bool {
	return len(strings.SplitAfter(message, fmt.Sprintf("/%s", command))) > 1
}

func getContentFromCommand(message string, command string) string {
	commands := strings.SplitAfter(message, fmt.Sprintf("/%s", command))
	if len(commands) > 1 {
		return strings.TrimSpace(commands[1])
	}
	return ""
}

// Builds and returns commands
func getCommands() []Command {

	return []Command{

		BotCommand{
			name:        "Help",
			description: "list of commands",
			matcher: func(update Update) bool {
				return update.Message.Text == "/help"
			},
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				var returnMsg bytes.Buffer
				returnMsg.WriteString("COMMANDS\n ")
				for _, command := range bot.GetCommands() {
					returnMsg.WriteString(fmt.Sprintf("\n<b>%s</b> - %s\n", command.Name(), command.Description()))
				}
				respChan <- *NewTextBotResponse(returnMsg.String(), update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "traps",
			description: "just a friendly reminder",
			matcher:     messageContainsCommandMatcher("traps"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("https://www.youtube.com/watch?v=Hmi3x-Ar86Q", update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "ping",
			description: "check if the bot is listening",
			matcher:     messageContainsCommandMatcher("ping"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("fuck you want?", update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "wuh",
			description: "stop fucking stop",
			matcher:     messageContainsCommandMatcher("wuh"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("https://www.youtube.com/watch?v=j3z7vjs1E18", update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "Ahem",
			description: "You are a furry arn't you",
			matcher: func(update Update) bool {
				return update.Message.Text == "ahem" && update.Message.ReplyToMessage != nil
			},
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse(fmt.Sprintf("%s is actually the furry", update.Message.ReplyToMessage.From.UserName), update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "Alias Set",
			description: "Alias the chat for other commands like edge, <code>/alias-set resistance</code>",
			matcher:     messageContainsCommandMatcher("alias-set"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				alias := getContentFromCommand(update.Message.Text, "alias-set")

				if alias != "" {
					_, alreadyExists := bot.IsAliasSet(alias)
					if alreadyExists {
						respChan <- *NewTextBotResponse(fmt.Sprintf("Someone has already taken the alias '%s'", alias), update.Message.Chat.ID)
					} else {
						bot.SetChatAlias(alias, update.Message.Chat.ID)
						respChan <- *NewTextBotResponse(fmt.Sprintf("Alias set as: '%s'", alias), update.Message.Chat.ID)
					}
				}
			},
		},

		BotCommand{
			name:        "Password",
			description: "Gives you chat id for edged site, /password",
			matcher:     messageContainsCommandMatcher("password"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse(strconv.FormatInt(update.Message.Chat.ID, 10), update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "Leaving",
			description: "Cause you want more attention, /leaving",
			matcher:     messageContainsCommandMatcher("leaving"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("Y’all are miserable people who demand to be right at all times, even when you have no experience on the subject. I’m sick of it and it’s gonna be better to just not have to deal with it. So I’m out. You’re mostly all hugely negative impacts on mine and others’ lives. Obviously some of you I still consider friends, but it’s really gotten to where the animosity of a few people make this little group irredeemably shitty for me to be a part of. Or irredeemably shitty in general really. There’s really just no good part of Resistance. ", update.Message.Chat.ID)
				respChan <- *NewTextBotResponse("Especially when it’s pretty much exclusively me that seems to be the target of all the hate. It just feels super mean spirited. It would be different if you attacked everybody the same way, but you don’t. And don’t even try to fucking pretend you do. I’ve genuinely never felt liked here. ", update.Message.Chat.ID)
			},
		},

                BotCommand{
			name:        "Opinion",
			description: "It's not what you said, it's how you said it",
			matcher:     messageContainsCommandMatcher("opinion"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("its not about your opinion. its just that you wont admit when you are wrong/agree with anothers view point OR you just scream insults", update.Message.Chat.ID)
				respChan <- *NewTextBotResponse("and no, i could give two fucks about your opinion. its your attitude i have a problem with.", update.Message.Chat.ID)
				respChan <- *NewTextBotResponse("the fact that you have reduced my comments calling you out on your shit attitude and bullshit to a \"rant\" just solidifies my point on how you view other people.", update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "God is great",
			description: "Are we a bad person",
			matcher: func(update Update) bool {
				return update.Message.Text == "/gg"
			},
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("GOD IS GREAT", update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "Edge",
			description: "Hide messages for later, reply to a message with /edge",
			matcher: func(update Update) bool {
				return update.Message.ReplyToMessage != nil && update.Message.Text == "/edge"
			},
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				bot.PushMessageToChatBuffer(strconv.FormatInt(update.Message.Chat.ID, 10), *update.Message.ReplyToMessage)

				whatWasEdged := "unknown message type. You're not getting that back lol"

				if update.Message.ReplyToMessage.Text != "" {
					whatWasEdged = "text.. the fuck could have been said to warrant edging text?"
				}

				log.Printf("%+v\n", update.Message.ReplyToMessage)
				if update.Message.ReplyToMessage.Photo != nil {
					photos := *update.Message.ReplyToMessage.Photo
					whatWasEdged = "photo"
					go bot.GetFile(photos[0].FileID)
				}

				if update.Message.ReplyToMessage.Sticker != nil {
					whatWasEdged = "sticker"
					sticker := *update.Message.ReplyToMessage.Sticker
					go bot.GetFile(sticker.FileID)
				}

				if update.Message.ReplyToMessage.Document != nil {
					whatWasEdged = "document"
				}

				if update.Message.ReplyToMessage.Video != nil {
					whatWasEdged = "video"
				}

				go bot.deleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
				go bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
				respChan <- *NewTextBotResponse(fmt.Sprintf("%s edged %s's %s", update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName, whatWasEdged), update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "Ejaculate",
			description: "Release all the messages that have been edged with /ejaculate",
			matcher: func(update Update) bool {
				return update.Message.Text == "/ejaculate"
			},
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				msgSentCount := 0
				buffer := bot.ClearBuffer(update.Message.Chat.ID)
				for msg := range buffer.Everything() {
					msgSentCount++
					respChan <- NewTextBotResponse(msg.From.UserName+" sent:", update.Message.Chat.ID)
					if msg.Photo != nil {
						photos := *msg.Photo
						respChan <- NewPictureReferenceBotResponse(photos[0].FileID, update.Message.Chat.ID)
					} else if msg.Sticker != nil {
						respChan <- NewStickerBotResponse(msg.Sticker.FileID, update.Message.Chat.ID)
					} else if msg.Document != nil {
						respChan <- NewFileReferenceBotResponse(msg.Document.FileID, update.Message.Chat.ID)
					} else if msg.Video != nil {
						respChan <- NewVideoReferenceBotResponse(msg.Video.FileID, update.Message.Chat.ID)
					} else if msg.Text != "" {
						respChan <- NewTextBotResponse(msg.Text, update.Message.Chat.ID)
					}
				}

				if msgSentCount == 0 {
					respChan <- *NewTextBotResponse("I'm not usually like this. Maybe if you do something sexy it'll start working", update.Message.Chat.ID)
				} else if msgSentCount < 5 {
					respChan <- *NewTextBotResponse("Normally I'm not that quick", update.Message.Chat.ID)
				} else if msgSentCount < 10 {
					respChan <- *NewTextBotResponse("I need a ciggarette after that", update.Message.Chat.ID)
				} else {
					respChan <- *NewTextBotResponse("HOLY FUCK I NEEDED THAT, sorry about the mess", update.Message.Chat.ID)
				}
			},
		},

		BotCommand{
			name:        "Rush",
			description: "Obliterate your opponent",
			matcher:     messageContainsCommandMatcher("rush"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				wholeCommand := getContentFromCommand(update.Message.Text, "rush")

				if wholeCommand == "" {
					return
				}
				commands := strings.Split(wholeCommand, " with ")

				name := commands[0]
				attack := ""

				if len(commands) > 1 {
					attack = commands[1]
				}
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					bot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 70,
				})

				for i := 0; i < 9; i++ {
					im, err := gg.LoadPNG(fmt.Sprintf("trunks2/F_00%d.png", i))
					if err != nil {
						bot.errorReport.Log("unable to load image: " + err.Error())
						return
					}
					dc := gg.NewContextForImage(im)
					dc.SetRGB(1, 0, 0)
					dc.SetFontFace(face)
					dc.DrawStringAnchored(update.Message.From.UserName, 950, 120, 0.0, 0.0)
					dc.DrawStringAnchored(name, 750, 600, 0.0, 0.0)
					dc.SavePNG(fmt.Sprintf("trunks2out/F_00%d.png", i))
				}

				if attack != "" {
					for i := 34; i < 42; i++ {
						im, err := gg.LoadPNG(fmt.Sprintf("trunks2/F_0%d.png", i))
						if err != nil {
							bot.errorReport.Log("unable to load image: " + err.Error())
							return
						}
						dc := gg.NewContextForImage(im)
						dc.SetRGB(1, 0, 0)
						dc.SetFontFace(face)
						dc.DrawStringAnchored(attack, 250, 300, 0.0, 0.0)
						dc.SavePNG(fmt.Sprintf("trunks2out/F_0%d.png", i))
					}
				}

				StichPicturesTogether("trunks2out")
				respChan <- *NewFileUploadBotResponse("movie.mp4", update.Message.Chat.ID)
			},
		},

		BotCommand{
			name:        "Bash",
			description: "Bully your opponent",
			matcher:     messageContainsCommandMatcher("bash"),
			execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				wholeCommand := getContentFromCommand(update.Message.Text, "bash")

				if wholeCommand == "" {
					return
				}

				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					bot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 70,
				})

				for i := 0; i < 9; i++ {
					im, err := gg.LoadPNG(fmt.Sprintf("bash/source/F_00%d.png", i))
					if err != nil {
						bot.errorReport.Log("unable to load image: " + err.Error())
						return
					}
					dc := gg.NewContextForImage(im)
					dc.SetRGB(1, 0, 0)
					dc.SetFontFace(face)
					dc.DrawStringAnchored(update.Message.From.UserName, 400, 100, 0.0, 0.0)
					dc.DrawStringAnchored(wholeCommand, 100, 100, 0.0, 0.0)
					dc.SavePNG(fmt.Sprintf("bash/out/F_00%d.png", i))
				}

				StichPicturesTogether("bash/out")
				respChan <- *NewFileUploadBotResponse("movie.mp4", update.Message.Chat.ID)
			},
		},

		holyCommand,
		killCommand,
		rule34Command,
		hedgehogCommand,
		saveCommand,
		pokedexCommand,
		swallowCommand,
		doitCommand,
		youwontCommand,
		repostCommand,
		mockCommand,
		yeetCommand,
		hmCommand,
		defineCommand,
		userIDCommand,
		stockPrice,
		resistanceRuleOneCommand,
		resistanceRuleTwoCommand,
		resistanceRuleThreeCommand,
		sarahSaysCommand,
	}

}

// Create our routes
func initRoutes(router *gin.Engine, telebot TeleBot) {
	router.SetFuncMap(template.FuncMap{
		"pictureDeref": func(i *[]PhotoSize) PhotoSize {
			if i == nil {
				return PhotoSize{}
			}

			photos := *i
			return photos[0]
		},
		"stickerDeref": func(i *Sticker) Sticker {
			if i == nil {
				return Sticker{}
			}
			return *i
		},
	})

	router.LoadHTMLGlob("templates/*.tmpl")

	timeStarted := GetTime()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"restarted": timeStarted,
			"errors":    telebot.errorReport.Generate(),
		})
	})

	router.GET("/edge/:chatID", func(c *gin.Context) {
		chatID := c.Param("chatID")
		msgs := telebot.ChatBuffer(chatID)
		c.HTML(http.StatusOK, "edge.tmpl", gin.H{
			"messages": msgs,
		})
	})

	router.StaticFS("/media", http.Dir("media"))
}

func listenForUpdates(telebot TeleBot) {

	for {
		// Sleep first, so if we error out and continue to the next loop, we still end up waiting
		time.Sleep(time.Second)

		updates, err := telebot.GetUpdates()

		if err != nil {
			telebot.errorReport.Log("Error getting updates from telegram: " + err.Error())
			continue
		}

		// Dispatch incoming messages to appropriate functions
		for _, update := range updates {
			if update.Message != nil {
				log.Println(update.Message.ToString())
				telebot.OnMessage(update)
			}
		}

	}
}

func logginToReddit(errorReport Report) RedditAccount {

	log.Printf("Logging into: %s\n", os.Getenv("REDDIT_USERNAME"))
	user, err := LoginToReddit(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		"Resistance Telegram Botter",
	)
	if err != nil {
		errorReport.Log("Error logging into reddit! " + err.Error())
	} else {
		log.Println(fmt.Sprintf("Succesfully logged in."))
	}

	return user
}

func main() {

	// Can't run a server without a port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable was not set")
		return
	}
	log.Printf("Starting bot using port %s\n", port)

	errorReport := NewReport()
	redditUser := logginToReddit(*errorReport)
	teleBot := NewTelegramBot(os.Getenv("TELE_KEY"), *errorReport, redditUser, getCommands())
	teleBot.Start()

	go listenForUpdates(*teleBot)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r, *teleBot)
	r.Run(":" + port)

}
