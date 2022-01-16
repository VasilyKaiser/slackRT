# slackRT
 Slack remote terminal - execute commands on remote host using slack slash command

 ## Installation

1. Go to api.slack.com/apps and sign in and create new app
2. Basic Information -> Add features and functionality -> Slash Commands -> create command `/sh` Request URL must match public URL/IP and of your remote host, port(you will need to enter it in config.yaml file).
3. In permissions add scope `chat:write` and `channels:history`.
4. Install app to your workspace.
5. In channel you want to receive app messages, Integrations -> Add apps -> add your app.
6. Fill in `config.yaml` file: `OAUTH_TOKEN` (from app OAuth & Permissions menu), `CHANNEL_ID` (channel you added your app to), `SIGNING_SECRET` (Basic Information -> App Credentials), `PORT` (Port must match with URL in slash command). Place this file in the same directory as slackRT.
7. Run slackRT on your remote host, or `go run main.go` from source.

## Usage

- `/sh ls -l` (or some other command available on your remote host)
- `/sh del` (this will remove messages from channel your bot has been writing to)