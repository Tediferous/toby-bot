# Toby-boT

## How to dev

### Requirements
1. go (atleast 1.14)
2. (optional) docker
3. a discord server you own

### Setup
0. Make a test bot account: https://discordpy.readthedocs.io/en/latest/discord.html
1. Invite bot to server
2. Take note of the bot account's token and keep it private
3. Enable developer Mode: https://discordia.me/en/developer-mode#:~:text=Enabling%20Developer%20Mode%20is%20easy,the%20toggle%20to%20enable%20it.
4. With dev mode enabled you can right click on various objects and copy their discord ID. Get the ID's of your:
  a. Spam Channel - I channel you would like to spam with messages
  b. Server 
  c. Ban Role - a role that only allows someone to speak in the banned channel
5. Install go modules with `go mod download`, i think?, maybe `go build`
6. run with `go run main.go -t <TOKEN>`

## How to add THE Toby to a Server
1. https://discord.com/oauth2/authorize?client_id=731318242903064667&scope=bot&permissions=1543576640

## Flags

t - Bot Token
s - Spam Channel ID for debug messages
b - Ban Role ID for ban function
g - Guild/Server ID

Example run command:
`toby -t <TOKEN> -s <SPAM ID> -b <BAN ROLE ID> -g <GUILD ID>`

## Commands

1. Ban - "@Toby ban @Ted, 2h", replaces all of users role with the ban role and returns them after time period


