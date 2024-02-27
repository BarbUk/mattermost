// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
)

const (
	InternalKeyPrefix = "mmi_"
	BotUserKey        = InternalKeyPrefix + "botid"
)

func initializePluginImplementation(pluginImplementation any) map[string]plugin.Plugin {
	if impl, ok := pluginImplementation.(interface {
		SetAPI(api API)
		SetDriver(driver Driver)
	}); !ok {
		panic("Plugin implementation given must embed plugin.MattermostPlugin")
	} else {
		impl.SetAPI(nil)
		impl.SetDriver(nil)
	}

	return map[string]plugin.Plugin{
		"hooks": &hooksPlugin{hooks: pluginImplementation},
	}
}

// Starts the serving of a Mattermost plugin over net/rpc. gRPC is not yet supported.
//
// Call this when your plugin is ready to start.
func ClientMain(pluginImplementation any) {
	pluginMap := initializePluginImplementation(pluginImplementation)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
	})
}

// Starts the serving of a Mattermost plugin over net/rpc in the context of a unit test.
//
// ctx is used to kill the plugin from the unit tests.
// reattachConfigCh receives the ReattachConfig to be sent back to the server.
// closeCh, if non-nil, will be closed when the plugin exits.
func ClientMainTesting(ctx context.Context, pluginImplementation any, reattachConfigCh chan<- *plugin.ReattachConfig, closeCh chan<- struct{}) {
	pluginMap := initializePluginImplementation(pluginImplementation)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
		Test: &plugin.ServeTestConfig{
			Context:          ctx,
			ReattachConfigCh: reattachConfigCh,
			CloseCh:          closeCh,
		},
	})
}

type MattermostPlugin struct {
	// API exposes the plugin api, and becomes available just prior to the OnActive hook.
	API    API
	Driver Driver
}

// SetAPI persists the given API interface to the plugin. It is invoked just prior to the
// OnActivate hook, exposing the API for use by the plugin.
func (p *MattermostPlugin) SetAPI(api API) {
	p.API = api
}

// SetDriver sets the RPC client implementation to talk with the server.
func (p *MattermostPlugin) SetDriver(driver Driver) {
	p.Driver = driver
}
