package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	db "github.com/sordid-rectangles/dev-tools-bot/revolver"
)

const Version = "v0.0.1-alpha"

var dg *discordgo.Session
var TOKEN string
var GUILDID string
var days = 1

type gun struct {
	chambers []bool
	loaded   bool
	bans     int
}

//var revolver = gun{chambers: []bool{false, false, false, false, false, false}, loaded: false, bans: 0}

func init() {
	// Print out a fancy logo!
	fmt.Printf(`Discord Roulette! %-16s\/`+"\n\n", Version)

	//Load dotenv file from .
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	//Load Token from env (simulated with godotenv)
	TOKEN = os.Getenv("BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Error loading token from env")
		os.Exit(1)
	}

	GUILDID = os.Getenv("GUILD_ID")
	if GUILDID == "" {
		log.Println("No GuildID specified in env")
		GUILDID = "" //this effectively specifies command registration as global
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "check",
			Description: "Checks if revolver is loaded",
		},
		{
			Name:        "load",
			Description: "Loads 1 revolver chamber",
		},
		{
			Name:        "safe",
			Description: "Unloads revolver",
		},
		{
			Name:        "spin",
			Description: "Spins the six chambers",
		},
		{
			Name:        "shoot",
			Description: "Pulls the trigger",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"check": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

			var guild = i.GuildID

			rev, exists := db.Memstore[guild]
			if exists {
				if rev.Loaded {
					content = "Revolver Loaded"
				} else {
					content = "Revolver Empty"
				}
			} else {
				db.Memstore[guild] = &db.Gun{GuildID: guild, Chambers: []bool{false, false, false, false, false, false}, Loaded: false, Bans: 0}
				content = "Revolver Empty"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"load": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

			var guild = i.GuildID

			rev, exists := db.Memstore[guild]
			if exists {
				rev.Load()
				content = "*Click!*"

			} else {
				db.Memstore[guild] = &db.Gun{GuildID: guild, Chambers: []bool{true, false, false, false, false, false}, Loaded: true, Bans: 0}
				content = "*Click!*"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"safe": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

			var guild = i.GuildID

			rev, exists := db.Memstore[guild]
			if exists {
				rev.Safe()
				content = "*Clink!*"

			} else {
				db.Memstore[guild] = &db.Gun{GuildID: guild, Chambers: []bool{false, false, false, false, false, false}, Loaded: true, Bans: 0}
				content = "*Clink!*"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"spin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

			var guild = i.GuildID

			rev, exists := db.Memstore[guild]

			var spun = false
			if exists {
				log.Println("Pre-spin")
				log.Println(rev.Chambers)
				spun = rev.Spin()

			} else {
				db.Memstore[guild] = &db.Gun{GuildID: guild, Chambers: []bool{true, false, false, false, false, false}, Loaded: true, Bans: 0}
				spun = false
			}

			if spun {
				log.Println("spun")
				log.Println(rev.Chambers)
				content = "*brrrrrrrrrr. Click*"
			} else {
				content = "You have to load the revolver first silly. We dont dry spin around here."
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"shoot": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

			var mem = i.Member

			nick := mem.Nick

			if nick == "" {
				nick = mem.User.Username
			}

			var guild = i.GuildID

			rev, exists := db.Memstore[guild]

			var fired = false
			if exists {
				log.Println("Pre-shoot")
				log.Println(rev.Chambers)
				fired = rev.Shoot()
				log.Println(fired)
				log.Println("Post-shoot")
				log.Println(rev.Chambers)

			} else {
				db.Memstore[guild] = &db.Gun{GuildID: guild, Chambers: []bool{true, false, false, false, false, false}, Loaded: true, Bans: 0}
				fired = false
			}

			if fired {
				content = fmt.Sprintf("BANG! \nGuess it wasn't %s's lucky day", nick)
				// reason := content
				// err := s.GuildBanCreateWithReason(i.GuildID, mem.User.ID, reason, days)
				t := time.Now().Add(time.Hour * 2)
				err := s.GuildMemberTimeout(i.GuildID, mem.User.ID, &t)
				if err != nil {
					log.Printf("Failed to timeout player. error: %s\nDuration of timeout: %v", err, time.Until(t).Hours())
				}

			} else {
				content = "CLICK! Lucky devil"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
	}
)

func init() {
	var err error
	dg, err = discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("Error creating discordgo session!")
		os.Exit(1)
	}
}

func main() {
	var err error
	//Configure discordgo session bot
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//Register Bot Intents with Discord
	//worth noting MakeIntent is a no-op, but I want it there for doing something with pointers later
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GUILDID, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)

	defer dg.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")

	// Exit Normally.
	//exit
}

func comesFromDM(s *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		if channel, err = s.Channel(i.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}
