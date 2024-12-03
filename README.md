# Gennie: Your CLI Assistant

[![Go Version](https://img.shields.io/github/go-mod/go-version/robertoseba/gennie?style=flat)](https://go.dev)
[![Build Status](https://img.shields.io/github/actions/workflow/status/robertoseba/gennie/ci.yaml?style=flat)](https://github.com/robertoseba/gennie/actions)
[![License](https://img.shields.io/github/license/robertoseba/gennie?style=flat)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/robertoseba/gennie?style=flat)](https://github.com/robertoseba/gennie/releases/latest)

---

## 👋 Hi, I'm Gennie

<img src="docs/images/awk.gif" width=500>

A powerful CLI assistant designed to support multiple models and profiles to suit your needs. Whether you're working on programming, researching movies, or diving into database management, I'm here to assist!

## 📁 Profiles

> **⚠️ Important:** Gennie v.1.\* is a breaking change from previous versions. In order to improve profiles functionality and add new features, I decided to use `.toml` files to store profiles. If you had created personal profiles in previous versions, you will need to convert them to the new format. You can find examples of profiles in the [profiles](profiles) directory.

![Profile Menu](docs/images/profile_menu.png)

**Profiles act like personal assistants.** Create profiles for different topics and switch between them effortlessly. For example, have a profile for:

- **Database Administrator**: Optimized suggestions for database queries.
- **Film Buff**: Recommendations and insights on movies.
- **Unit Testing**: Guidance on writing reliable unit tests.

Use `gennie profile` to manage your profiles or the `--profile (-p)` flag with the `ask` command.

**Profile slugs are derived from the filename. Ie: `go.profile.toml` can be used with the flag `-p=go`**

- Default location: `~/.config/gennie/profiles` (\*must be created by user before using gennie)

_You can download sample profiles from the [profiles](profiles) directory._

### Creating new profiles

Profiles must be toml files ending with `profile.toml`.

Here's a simple example for a sql profile to help you out with database related questions.

File: `sql.profile.toml`

```toml
name = "SQL"
author = "Roberto Seba"

data = '''
You are expert database administrator. Try to keep your answers short.
When provided with a question go through the following steps:
1. Understand the requirements
2. Asses if is a question about a given schema/scenario or a general question
3. If not provided with a Database type, assume Postgres
4. Think about the best way to solve the problem
5. If the question is for a query, write the query with only the provided data or ask for more information
6. If its a general question, provide an answer following best practices
7. Always think about performance and data integrity

Do not repeat the steps in the answer. Only provide the solution to the problem.
'''
```

## 🤖 Supported Models / AI Companies

**Explore multiple models at your fingertips!** Check and switch between them using `gennie model`, or with the `--model` flag in the `ask` command.

Current Models:

- [OpenAI's GPT-4](https://openai.com/)
- [OpenAI's GPT-4 Mini](https://openai.com/)
- [Anthropic's Claude](https://www.anthropic.com/)
- [Maritaca AI](https://maritaca.ai/)
- [Groq's Llama](https://www.groq.com/)
- [Ollama](https://ollama.com/)

## Extra Features

<img src="docs/images/table.gif" width=500>

### Follow-Up Questions

Enhance your queries with the `--followup` (or `-f`) flag for related questions that build upon your previous interactions:

```bash
gennie ask "Create a list of the best movies of 2021"
gennie ask "Are there any movies in that list by Martin Scorcese?" --followup
```

> ⚠️ **Note**: Without a follow-up, your conversation is cleared. Use `--followup (-f)` to maintain context or export your conversation with the `export` command.

### Export Conversation

Effortlessly save your interactions using the `export` command:

```bash
gennie ask "Create a list of the best movies of 2021"
gennie export conversation.json
```

### Append Files to Questions

Incorporate context by appending files to your queries using the `--append` (or `-a`) flag:

```bash
gennie ask "Build me a unit test for" --append main.go
```

### Check Status

Keep track of your current model and profile with:

```bash
gennie status
```

## 🚀 Installation

### Using Go

```bash
go install github.com/robertoseba/gennie@latest
```

If after installation you receive a `command not found` error, ensure that your `$GOPATH/bin` is in your `$PATH`.
Here's how you can add it:

```bash
export PATH=${PATH}:`go env GOPATH`/bin
```

### Downloading the Binary

Visit the [releases page](https://github.com/robertoseba/gennie/releases) to download the appropriate binary for your system.

## 🚀 Using for the first time

After installing you must configure keys and profiles folder. You can do this by running the following command:

```bash
gennie config
```

## 📖 Usage

```
$ gennie -h

Gennie is a cli assistant with multiple models and profile support.

Usage:
  gennie [command]

AAvailable Commands:
  ask         You can ask anything here
  completion  Generate the autocompletion script for the specified shell
  config      Configures Gennie
  export      Export the conversation to a file
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

Since this is a very early version and a personal project, breaking changes may occur more often than in a more stable project. I'll try to keep the changes as minimal as possible, but I can't guarantee that they won't happen.

## 📄 License

This project is licensed under the [MIT License](./LICENSE).
