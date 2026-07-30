// BetterDiscord's Injection Script (app.asar method)
const path = require("path");
const electron = require("electron");

// The global BetterDiscord folder lives one directory above userData (the
// appData root). Electron gives the postfixed userData, so go up a directory.
let userConfig = path.join(electron.app.getPath("userData"), "..");

// If we're on Linux there are a couple cases to deal with
if (process.platform !== "win32" && process.platform !== "darwin") {
    // Use || instead of ?? because a falsey value of "" is invalid per XDG spec
    userConfig = process.env.XDG_CONFIG_HOME || path.join(process.env.HOME, ".config");
}

require(path.join(userConfig, "BetterDiscord", "data", "betterdiscord.asar"));

// Hand off to Discord's real (renamed) app entry point
module.exports = require("../betterdiscord.app.asar");