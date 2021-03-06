# Tech Talks Twitch bot

This is a simply Twitch bot oriented to provide useful features for Tech Talks


## Getting started

Copy `.env.example` to `.env` and update the values

#### `BOT_USERNAME` 
Bot user name

#### `TMI_OAUTH_TOKEN`
Twitch Oauth token. Can generate it on https://twitchapps.com/tmi/
Be sure the user logged to generate this token is the channel owner in order to have higher chat publish rate.

#### `CHANNEL`
Channel where the bot will publish and observe

#### `COMMAND_COLD_DOWN_TIME`
Time in seconds to cold down every command, in this time the bot only will publish once for every command

#### `COMMAND_HELP_EVERY`
Time in second for sending the help message to chat. **0 disables the automatic command **

#### `DISABLED_COMMANDS`
Ids of the disabled commands (separated by commas)

#### `DISABLE_CHAT_LOGGER`
If true, disables the chat logger

### Configure the commands

Copy or rename all `data/*.json.example` to `data/*.json` and edit this files to configure every command


## Running

Build and then run

### Linux / MacOs
```bash
$ go build cmd/main.go -o bot
$ ./bot
# or
$ go run cmd/main.go
```

### Windows
```bash
$ go build -o bot.exe
$ bot.exe
```


## Features
The bot will say to the chat on start and every `COMMAND_HELP_EVERY`s the list of the available commands
If its 0 no help command
Every command has its own cold down time to avoid flooding the chat 


## Current commands

* `!ponentes` Say the list of meetup/event speakers
* `!videos` Say the last video, and the youtube channel from Vigotech.json or /data/videos.json
* `!eventos` Say upcomming events Vigotech.json or /data/events.json

### Text command
Editing `data/text.json` you can define multiple commands that responds with a simple text:
```json
[
  {
    "command": "descripcion",
    "text": "2021 es un año que tiene ...."
  }
]
```
for example:
    * `!descripcion` Say the meetup/event description


## Add new command

Create a new file in `commands/`
Command must return a pointer to struct

```go
type Command struct {
	Id string // Unique id
	Name string // Command
	client *twitch.Client // Pointer to client
	handler func(client *twitch.Client, message twitch.PrivateMessage) error
}
```
The easy way to do that is create and "Constructor"


```go
func NewYourCommand(client *twitch.Client) *Command {
    return &Command{
        Id:     "your-command-id",
        Name:   "your-command",
        client: client,
        handler: func(client *twitch.Client, message twitch.PrivateMessage) error {
            // Do things
            return nil
        },
    }
}
```



In the file `commands/commands.go` add the new command struct to the function get commands
```go
func GetCommands (client *twitch.Client) []*Command {
    return []*Command{
        NewDescriptionCommand(client),
        NewSpeakersCommand(client),
        NewYourCommand(client),
        ...,
    }
}
```



