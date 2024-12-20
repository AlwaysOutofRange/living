package main

import (
	"fmt"
	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sirupsen/logrus"
	"log/slog"
	"time"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := server.DefaultConfig().Config(slog.Default())
	if err != nil {
		log.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for p := range srv.Accept() {
		accept(p)
	}
}

func accept(p *player.Player) {
	opts := world.EntitySpawnOpts{
		Position: p.Position(),
	}

	conf := living.Config{
		EntityType: entityTypeEnderman{},
		Handler:    handler{},
		MaxHealth:  40,
		Drops: []living.Drop{
			living.NewDrop(item.EnderPearl{}, 0, 2),
		},
	}
	p.Tx().AddEntity(opts.New(conf.EntityType, conf))
}

type entityTypeEnderman struct {
	living.NopLivingType
}

func (entityTypeEnderman) EncodeEntity() string {
	return "minecraft:enderman"
}
func (entityTypeEnderman) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 2.9, 0.3)
}

type handler struct {
	living.NopHandler
}

func (handler) HandleHurt(ctx *living.Context, damage float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	fmt.Println("enderman hurt")
}
