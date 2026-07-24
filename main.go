package main

import (
	"github.com/huntbase-io/steampipe-plugin-anthropic/anthropic"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: anthropic.Plugin,
	})
}
