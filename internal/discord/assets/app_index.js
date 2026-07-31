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

// Never let a missing or broken BetterDiscord asar keep Discord from launching:
// this file is the app entry point, so an unhandled throw here bricks the client.
// Log and fall through to Discord's real app.
try {
    require(path.join(userConfig, "BetterDiscord", "data", "betterdiscord.asar"));
}
catch (error) {
    console.error("Failed to load BetterDiscord:", error);
}

// Hand off to Discord's real (renamed) app entry point
module.exports = require("../betterdiscord.app.asar");