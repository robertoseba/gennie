package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
)

//go:embed version.txt
var version string

func main() {
	cmdUtil, err := cmd.NewCmdUtil(version)
	if err != nil {
		cmd.ExitWithError(err)
	}
	defer cmdUtil.Storage.Save()

	command := cmd.NewRootCmd(cmdUtil)
	if cmdUtil.Storage.IsNew() {
		command.SetArgs([]string{"config"})
	}

	command.Execute()
}

//TODO:
// Main should loadCache(), cache.GetConfig(), cache.GetProfilesInfo(), cache.GetActiveSession()
// session: {profileSlug, modelSlug, conversation}
// session.SetModel(model), session.SetProfile(profile)
// NewProfileRepository(cache.GetProfilesInfo())
//cmd.NewRootCmd(config, cache)
// cache.setConfig(), cache.setProfilesInfo(), cache.setActive()
// cache.save()
