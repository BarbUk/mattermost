// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
package app

import (
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

// ReattachPlugin allows the server to bind to an existing plugin instance launched elsewhere.
func (a *App) ReattachPlugin(manifest *model.Manifest, pluginReattachConfig *model.PluginReattachConfig) *model.AppError {
	return a.ch.ReattachPlugin(manifest, pluginReattachConfig)
}

// ReattachPlugin allows the server to bind to an existing plugin instance launched elsewhere.
func (ch *Channels) ReattachPlugin(manifest *model.Manifest, pluginReattachConfig *model.PluginReattachConfig) *model.AppError {
	pluginsEnvironment := ch.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return model.NewAppError("ReattachPlugin", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	// Deactivate and remove any existing plugin, if present.
	pluginsEnvironment.Deactivate(manifest.Id)
	pluginsEnvironment.RemovePlugin(manifest.Id)

	// Reattach to the plugin
	if err := pluginsEnvironment.Reattach(manifest, pluginReattachConfig); err != nil {
		return model.NewAppError("ReattachPlugin", "app.plugin.reattach.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return nil
}
