# Gennie: Your CLI Assistant

[![Go Version](https://img.shields.io/github/go-mod/go-version/robertoseba/gennie?style=flat)](https://go.dev)
[![Build Status](https://img.shields.io/github/actions/workflow/status/robertoseba/gennie/ci.yaml?style=flat)](https://github.com/robertoseba/gennie/actions)
[![License](https://img.shields.io/github/license/robertoseba/gennie?style=flat)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/robertoseba/gennie?style=flat)](https://github.com/robertoseba/gennie/releases/latest)

---

## 👋 Hi, I'm Gennie!

<img src="docs/images/awk.gif" width=500>

A powerful CLI assistant designed to support multiple models and profiles to suit your needs. Whether you're working on programming, researching movies, or diving into database management, I'm here to assist!

## 📁 Profiles

![Profile Menu](docs/images/profile_menu.png)

**Profiles act like personal assistants.** Create profiles for different topics and switch between them effortlessly. For example, have a profile for:

- **Database Administrator**: Optimized suggestions for database queries.
- **Film Buff**: Recommendations and insights on movies.
- **Unit Testing**: Guidance on writing reliable unit tests.

Use `gennie profile` to manage your profiles or the `--profile` flag with the `ask` command.

**Profiles are cached locally for performance:**

- Default location: `~/.config/gennie/profiles`
- Refresh your cached profiles with `gennie profile refresh`.

_You can download sample profiles from the [profiles](profiles) directory._

### Creating new profiles

Profiles must be json files ending with `profile.json`.

Here's a simple example for a sql profile to help you out with database related questions.

File: `sql.profile.json`

```json
{
  "name": "SQL", //This is the name that will show up in the profile menu. Can be more descriptive than the slug
  "slug": "sql", //Slug is used to identify the profile when using the --profile flag
  "author": "Roberto Seba",
  "data": "You are expert database administrator especially in MySQL and PostgreSQL. Try to keep your answers short. It's always important to think about query performance and data integrity." //Data is where you prep the assistant before you ask questions. It can be as long as you want.
}
```

## 🤖 Supported Models / AI Companies

**Explore multiple models at your fingertips!** Check and switch between them using `gennie model`, or with the `--model` flag in the `ask` command.

Current Models:

- [OpenAI's GPT-4](https://openai.com/)
- [OpenAI's GPT-4 Mini](https://openai.com/)
- [Anthropic's Claude](https://www.anthropic.com/)
- [Maritaca AI](https://maritaca.ai/)
- [Groq's Llama](https://www.groq.com/)

### Coming Soon:

- [Ollama](https://ollama.com/)

## Extra Features

<img src="docs/images/table.gif" width=500>

### Follow-Up Questions

Enhance your queries with the `--followup` flag for related questions that build upon your previous interactions:

```bash
$ gennie ask "Create a list of the best movies of 2021"
$ gennie ask "Are there any movies in that list by Martin Scorcese?" --followup
```

> ⚠️ **Note**: Without a follow-up, your chat history is cleared. Use `--followup` to maintain context or export your history with the `export` command.

### Export Chat History

Effortlessly save your chat interactions using the `export` command:

```bash
$ gennie ask "Create a list of the best movies of 2021"
$ gennie export chat_history.txt
```

### Append Files to Questions

Incorporate context by appending files to your queries using the `--append` flag:

```bash
$ gennie ask "Build me a unit test for" --append main.go
```

### Check Status

Keep track of your current model and profile with:

```bash
$ gennie status
```

## 🚀 Installation

### Using Go:

```bash
$ go install github.com/robertoseba/gennie@latest
```

If after installation you receive a `command not found` error, ensure that your `$GOPATH/bin` is in your `$PATH`.
Here's how you can add it:

```bash
export PATH=${PATH}:`go env GOPATH`/bin
```

### Downloading the Binary:

Visit the [releases page](https://github.com/robertoseba/gennie/releases) to download the appropriate binary for your system.

## 🚀 Using for the first time

After installing you must configure keys and profiles folder. You can do this by running the following command:

```bash
$ gennie config
```

## 📖 Usage

```
$ gennie -h

Gennie is a cli assistant with multiple models and profile support.

Usage:
  gennie [command]

AAvailable Commands:
  ask         You can ask anything here
  clear       Clears all the conversation and preferences from cache
  completion  Generate the autocompletion script for the specified shell
  config      Configures Gennie
  export      Export the chat history to a file
  help        Help about any command
  model       Configures the model to use.
  profile     Profile management
  status      Shows the current status of gennie

Flags:
  -h, --help      help for gennie
  -v, --version   version for gennie

Use "gennie [command] --help" for more information about a command.
```

## API Keys

Each model requires an API key to function.
Use the `gennie config` command to set your API keys.

## 🐛 Issues and Suggestions

Gennie is an **OPEN** source project in its early stages. We welcome any bugs, issues, or suggestions you may have. Feel free to create an issue or contact me directly, and I'll respond as soon as possible!

## 📄 License

This project is licensed under the [MIT License](./LICENSE).
