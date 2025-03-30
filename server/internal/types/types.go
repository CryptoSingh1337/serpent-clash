package types

import "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"

type Id uint32

type Component interface {
	component.Input | component.Network | component.PlayerInfo | component.Snake
}
