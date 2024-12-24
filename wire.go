//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"instaLinkCollector/adapter"
	"instaLinkCollector/controller"
)

func InitVideoController() (*controller.VideoController, error) {
	wire.Build(
		controller.NewVideoController,
		adapter.NewVideoService,
	)
	return nil, nil
}
