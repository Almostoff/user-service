package utils

import (
	"github.com/goombaio/namegenerator"
	"github.com/mcbattirola/avatargen"
	"time"
)

func GenerateNickName() string {

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()

	return name
}

func GenerateAvatar(nickname string) string {
	svg := avatargen.Generate(nickname)
	return svg
}
