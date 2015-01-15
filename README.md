# Remote client for Tikplay

A commandline client for the TikPlay server that routes traffic through a ssh tunnel.

# Usage

The client is used with syntax: `tikp <command> <parameter>`
### Commands
* `play <url>` = Plays the given url
* `np` = Shows currently playing song
* `list <n>` = shows the n songs from the queue (n defaults to 10)
* `skip` = skips the currently playing song
* `clear` = clears the whole queue
