# PendaGo

A discord bot for PNH server in Golang, using [discorgo](https://github.com/bwmarrin/discordgo) !

# Build

Have Go environnement setup on your machine.

```
go build
```

# Usage

```
Usage of ./pendago:
  -c string
        JSON config file
  -t string
        Bot token
```

# RCON

The bot uses `mcrcon` to send Minecraft whitelist commands. Ensure it is installed
and available in `PATH` on the host running the bot.

# Config file

This is an example of the `config.json` required as `-c` parameter :

```json
{
  "wow_log_chan_id": "<CHAN ID>",
  "wow_chan_id": "<CHAN ID>",
  "bot_chan_id": "<CHAN ID>",
  "role_id": "<ROLE ID>",
  "trigger_command": "<CMD>",
  "role_boost": "<ROLE ID>",
  "penda_role": "<ROLE ID>",
  "penda_gold_role": "<ROLE ID>",
  "sw_command": "<CMD>",
  "swc_role_id": "<ROLE ID>",
  "rcon_host": "<RCON HOST>",
  "rcon_port": "<RCON PORT>",
  "rcon_password": "<RCON PASSWORD>"
}
```
