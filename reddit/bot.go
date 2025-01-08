package reddit

import (
	"net/http"
	"time"
)

// BotConfig configures a Reddit bot's behavior with the Reddit package.
type BotConfig struct {
	// Agent is the user-agent sent in all requests the bot makes through
	// this package.
	Agent string
	// App is the information for your registration on Reddit.
	// If you are not familiar with this, read:
	// https://github.com/reddit/reddit/wiki/OAuth2
	App App
	// Rate is the minimum amount of time between requests. If Rate is
	// configured lower than 1 second, the it will be ignored; Reddit's API
	// rules cap OAuth2 clients at 60 requests per minute. See package
	// overview for rate limit information.
	Rate time.Duration
	// Custom HTTP client
	Client *http.Client
}

// Bot defines the behaviors of a logged in Reddit bot.
type Bot interface {
	Account
	Lurker
	Scanner
	Reaper
}

type bot struct {
	Account
	Lurker
	Scanner
	Reaper
}

// NewBot returns a logged in handle to the Reddit API.
func NewBot(c BotConfig) (Bot, error) {
	cli, err := newClient(clientConfig{agent: c.Agent, app: c.App, client: c.Client})
	r := newReaper(
		reaperConfig{
			client:   cli,
			parser:   newParser(),
			hostname: "oauth.reddit.com",
			tls:      true,
			rate:     maxOf(c.Rate, time.Second),
		},
	)
	return &bot{
		Account: newAccount(r),
		Lurker:  newLurker(r),
		Scanner: newScanner(r),
		Reaper: r,
	}, err
}

// NewBotFromAgentFile calls NewBot with a config built from an agent file. An
// agent file is a convenient way to store your bot's account information. See
// https://github.com/turnage/graw/wiki/agent-files
func NewBotFromAgentFile(filename string, rate time.Duration) (Bot, error) {
	agent, app, err := load(filename)
	if err != nil {
		return nil, err
	}

	return NewBot(
		BotConfig{
			Agent: agent,
			App:   app,
			Rate:  rate,
		},
	)
}
