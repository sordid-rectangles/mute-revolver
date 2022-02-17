package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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

var revolver = gun{chambers: []bool{false, false, false, false, false, false}, loaded: false, bans: 0}

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
			Description: "Reloads revolver",
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

			if revolver.loaded {
				content = "Revolver Loaded"
			} else {
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

			revolver.load()
			content = "*Click!*"

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"safe": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			revolver.safe()
			content = "*Clink!*"

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content),
				},
			})
		},
		"spin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			log.Println("Pre-spin")
			log.Println(revolver.chambers)
			spun := revolver.spin()

			if spun {
				log.Println("spun")
				log.Println(revolver.chambers)
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
			var mem = i.Member
			nick := mem.Nick
			if nick == "" {
				nick = mem.User.Username
			}

			log.Println("Pre-shoot")
			log.Println(revolver.chambers)
			fired := revolver.shoot()
			log.Println(fired)
			log.Println("Post-shoot")
			log.Println(revolver.chambers)

			if fired {
				content = fmt.Sprintf("BANG! \nGuess it wasn't %s's lucky day", nick)
				reason := content
				err := s.GuildBanCreateWithReason(i.GuildID, mem.User.ID, reason, days)
				if err != nil {
					log.Printf("Failed to ban player. error: %s", err)
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

func (g *gun) shoot() bool {
	if g.loaded {
		if g.chambers[0] {
			g.loaded = false
			g.chambers[0] = false
			return true
		} else {
			newChambers := append(g.chambers[1:], false)
			g.chambers = newChambers
			return false
		}

	} else {
		return false
	}
}

func (g *gun) spin() bool {
	if g.loaded {
		g.chambers = []bool{false, false, false, false, false, false}
		rand.Seed(time.Now().UnixNano())
		l := len(g.chambers)
		i := rand.Intn(l)
		g.chambers[i] = true
		return true
	} else {
		return false
	}
}

func (g *gun) load() {
	g.chambers = []bool{true, false, false, false, false, false}
	g.loaded = true
}

func (g *gun) safe() {
	g.chambers = []bool{false, false, false, false, false, false}
	g.loaded = false
}
