# Remote client for Tikplay

A commandline client for the TikPlay server that routes traffic through a ssh tunnel.

# Installation

run `go install` in the tikp folder. If you get error messages about missing dependecies run `go get` for each of them.

The config.gcfg needs your ssh username, but the password is optional if you have ssh-agent running.

# Usage

The client is used with syntax: `tikp <command> <parameter>`
### Commands
* `play <url>` = Plays the given url
* `np` = Shows currently playing song
* `list <n>` = shows the n songs from the queue (n defaults to 10)
* `skip` = skips the currently playing song
* `clear` = clears the whole queue
* `task <id>` = shows the status of task with given id
