package revolver

import (
	"math/rand"
	"time"
)

type Gun struct {
	GuildID  string
	Chambers []bool
	Loaded   bool
	Bans     int
}

//Local map of GuildID => Gun{	guildID  string, chambers []bool, loaded   bool, bans     int}
var Memstore = make(map[string]*Gun)

func encode(*Gun) (bool, error) {
	var err error = nil
	return false, err
}

func (g *Gun) Load() {
	g.Chambers = []bool{true, false, false, false, false, false}
	g.Loaded = true
}

func (g *Gun) Shoot() bool {
	if g.Loaded {
		if g.Chambers[0] {
			g.Loaded = false
			g.Chambers[0] = false
			return true
		} else {
			newChambers := append(g.Chambers[1:], false)
			g.Chambers = newChambers
			return false
		}

	} else {
		return false
	}
}

func (g *Gun) Spin() bool {
	if g.Loaded {
		g.Chambers = []bool{false, false, false, false, false, false}
		rand.Seed(time.Now().UnixNano())
		l := len(g.Chambers)
		i := rand.Intn(l)
		g.Chambers[i] = true
		return true
	} else {
		return false
	}
}

func (g *Gun) Safe() {
	g.Chambers = []bool{false, false, false, false, false, false}
	g.Loaded = false
}
