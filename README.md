# 🎮 Recco

A terminal app that fetches your Steam library and uses **Gemini AI** to recommend games based on your selection.

Built with [charmbracelet/huh](https://github.com/charmbracelet/huh) for interactive TUI forms.

## Demo

![recco-demo-an](https://github.com/user-attachments/assets/3dca522e-31ba-48bf-91ab-c77e9c1eea92)


## Setup

1. **Steam API Key** → https://steamcommunity.com/dev/apikey
2. **Steam ID (64-bit)** → https://steamid.io
3. **Gemini API Key** → https://aistudio.google.com/app/apikey

Copy `.env.example` to `.env` and fill in your keys:

```sh
cp .env.example .env
# edit .env
```

## Run

```sh
go run .
```

Or build a binary:

```sh
go build -o steam-bubbles .
./steam-bubbles
```

## How it works

1. Fetches your Steam library via the `IPlayerService/GetOwnedGames` API
2. Displays your games sorted by playtime in a searchable multi-select list
3. Sends your selected games to Gemini for tailored recommendations
4. Pretty-prints the results in a styled terminal box
