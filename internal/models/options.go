package models

type InstallOptions struct {
	RestartDiscord bool `json:"restartDiscord"`
	UseDevBuild    bool `json:"useDevBuild"`
}

type RepairOptions struct {
	DisablePlugins       bool `json:"disablePlugins"`
	DisableThemes        bool `json:"disableThemes"`
	ClearCustomCSS       bool `json:"clearCustomCSS"`
	ClearWebpackCache    bool `json:"clearWebpackCache"`
	ClearAddonStoreCache bool `json:"clearAddonStoreCache"`
	ResetSettings        bool `json:"resetSettings"`
}

type UninstallOptions struct {
	FullUninstall  bool `json:"fullUninstall"`
	RestartDiscord bool `json:"restartDiscord"`
}
