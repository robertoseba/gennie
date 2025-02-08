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

- Default location: `~/.config/gennie/profiles` (\*must be created by user before using gennie).
  Every file inside the profile directory ending with `.profile.toml` will be automatically read as a profile

_You can download sample profiles from the [profiles](profiles) directory._

### Creating new profiles

Profiles must be toml files ending with `profile.toml`.

Here's a simple example, with the required fields, for a sql profile to help you out with database related questions.

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

### Using a profile directly in a question

```bash
gennie ask "how can I read a EXPLAIN returned from PostgreSQL" -p=sql
```

## 🤖 Supported Models / AI Companies

**Explore multiple models at your fingertips!** Check and switch between them using `gennie model`, or with the `--model` flag in the `ask` command.

Current Models:

- [OpenAI's GPT-4](https://openai.com/)
- [OpenAI's GPT-4 Mini](https://openai.com/)
- [Anthropic's Claude](https://www.anthropic.com/)
- [Maritaca AI](https://maritaca.ai/)
- [Groq's DeepSeek-R1-Distill-Llama-70B](https://www.groq.com/)
- [Ollama](https://ollama.com/)

### List models slugs for using directly in a question

```bash
gennie model slugs
```

### Using a model directly in a question

```
gennie ask who won the oscar for best movie in 2023 -m=sonnet
```

## Conversations

This is how you can manage your conversation history and keep your data for future use.
The last question is always saved as a active conversation. You can always use follow up to continue on the same conversation. BUT if you don't use the flag, a new conversation is started and you lose the previous one.
So, in order for you to maintain an archive of your conversations you can use these commands: `conversation save` and `conversation load`
That way you can have a directory with all you conversations for whenever you want to continue the chat.

## Saving the active conversation for future use

Let's say you made a bunch of questions but you feel like you might want to revisit this conversation.
Save it using:

```
gennie conversation save <jsonfilename>
```

### Reloading past conversation into active conversation

You can use this whenever you want to reload old conversations and make it active again

```
gennie conversation load <filename of previously saved conversation>
```

## Extra Features

<img src="docs/images/table.gif" width=500>

### Follow-Up Questions

Enhance your queries with the `--followup` (or `-f`) flag for related questions that build upon your previous interactions:

```bash
gennie ask "Create a list of the best movies of 2021"
gennie ask "Are there any movies in that list by Martin Scorcese?" --followup
```

> ⚠️ **Note**: Without a follow-up, your conversation is cleared. Use `--followup (-f)` to maintain context or export your conversation with the `export` command.

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

After installing you must configure api keys for the llm providers (OpenAi, Anthropic...) and a profiles folder where you'll keep your profile's collection. You can do this by running the following command:

```bash
gennie config
```

## 📖 Help

```
gennie --help
```

## API Keys

Each model requires an API key to function.
Use the `gennie config` command to set your API keys.

## 🐛 Issues and Suggestions

Gennie is an **OPEN** source project in its early stages. We welcome any bugs, issues, or suggestions you may have. Feel free to create an issue or contact me directly, and I'll respond as soon as possible!

Since this is a very early version and a personal project, breaking changes may occur more often than in a more stable project. I'll try to keep the changes as minimal as possible, but I can't guarantee that they won't happen.

## 📄 License

This project is licensed under the [MIT License](./LICENSE).
