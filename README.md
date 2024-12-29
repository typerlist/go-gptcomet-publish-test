# GPTComet
<!-- TOC -->

- [GPTComet](#gptcomet)
    - [Description](#description)
    - [See Also](#see-also)
    - [Features](#features)
    - [Installation](#installation)
        - [Prerequisites](#prerequisites)
        - [Build from Source](#build-from-source)
    - [Usage](#usage)
        - [Basic Usage](#basic-usage)
        - [Using a specific language](#using-a-specific-language)
        - [Dry Run](#dry-run)
        - [Configuring a New Provider](#configuring-a-new-provider)
        - [Managing Configuration](#managing-configuration)
        - [Generating Rich Commit Messages](#generating-rich-commit-messages)
    - [Configuration](#configuration)
        - [Supported Configuration Keys](#supported-configuration-keys)
    - [Contribution](#contribution)
    - [Testing](#testing)
    - [License](#license)
    - [Acknowledgements](#acknowledgements)

<!-- /TOC -->
## Description

GPTComet is an AI-powered command-line tool that generates and creates conventional Git commit messages based on staged changes. It leverages various Large Language Models (LLMs) to analyze the diff and produce meaningful commit messages, streamlining the development workflow.

## See Also

Python version: https://github.com/belingud/gptcomet

## Features

-   **Multiple LLM Providers:** Supports a wide range of LLM providers, including OpenAI, Claude, Gemini, Mistral, Ollama, and more.
-   **Interactive Provider Configuration:** Easily configure new providers and their settings through an interactive command-line interface.
-   **Customizable Commit Message Style:** Generate both brief and detailed commit messages using configurable prompt templates.
-   **Commit Message Translation:** Translate generated commit messages into various languages.
-   **Filtered Diff:** Ignores specified files (e.g., lock files, specific file types) in the diff to generate more relevant commit messages.
-   **Flexible Configuration:** Manage settings through a YAML configuration file, supporting options like API keys, models, output language, and more.
-   **Dry Run Mode:** Preview the generated commit message without actually creating a commit.
-   **Rich Commit Messages:**  Generate detailed commit messages with a title, summary, and bullet points describing the changes when using the `--rich` flag.

## Installation

### Prerequisites

-   Go (version 1.18 or later)

### Build from Source

1. Clone the repository:
    ```bash
    git clone https://github.com/belingud/go-gptcomet.git
    cd go-gptcomet
    ```
2. Build the project:
    ```bash
    go build -o gptcomet .
    ```
3. (Optional) Move the `gptcomet` binary to a directory in your `PATH` for easy access.

## Usage

### Basic Usage

1. Stage your changes in a Git repository:
    ```bash
    git add .
    ```
2. Run GPTComet to generate and create a commit message:
    ```bash
    ./gptcomet commit
    ```
    The tool will analyze the staged changes, generate a commit message, and prompt you to confirm or edit it before creating the commit.

### Using a specific language

You can change the output language by setting the `output.lang` configuration option:

```bash
./gptcomet config set output.lang zh-cn
```

### Dry Run

To preview the generated commit message without committing, use the `--dry-run` flag:

```bash
./gptcomet commit --dry-run
```

### Configuring a New Provider

To configure a new LLM provider:

```bash
./gptcomet newprovider
```

This will guide you through selecting a provider and entering the required configuration values (e.g., API key, model name).

### Managing Configuration

The `gptcomet config` command provides subcommands for managing the configuration file:

-   `get <key>`: Get the value of a configuration key.
-   `list`: List the entire configuration content.
-   `reset`: Reset the configuration to default values (optionally reset only the prompt section with `--prompt`).
-   `set <key> <value>`: Set a configuration value.
-   `path`: Get the configuration file path.
-   `remove <key> [value]`: Remove a configuration key or a value from a list.
-   `append <key> <value>`: Append a value to a list configuration.
-   `keys`: List all supported configuration keys.

**Example:**

```bash
./gptcomet config set openai.api_key "your_openai_api_key"
./gptcomet config get output.lang
./gptcomet config list
```

### Generating Rich Commit Messages

To generate a more detailed commit message, use the `--rich` flag:

```bash
./gptcomet commit --rich
```

This will use the `rich_commit_message` prompt template, resulting in a commit message with a title, summary, and potentially bullet points outlining the changes.

## Configuration

GPTComet stores its configuration in a YAML file located at `~/.config/gptcomet/gptcomet.yaml`.

### Supported Configuration Keys

Here's a summary of the main configuration keys:

| Key                             | Description                                                                                                  | Default Value            |
| :------------------------------ | :----------------------------------------------------------------------------------------------------------- | :----------------------- |
| `provider`                      | The name of the LLM provider to use.                                                                       | `openai`                 |
| `file_ignore`                   | A list of file patterns to ignore in the diff.                                                               | (See `config.go`)      |
| `output.lang`                   | The language for commit message generation.                                                                  | `en`                     |
| `output.rich_template`          | The template to use for rich commit messages.                                                              | `<title>:<summary>\n\n<detail>` |
| `console.verbose`               | Enable verbose output.                                                                                       | `true`                    |
| `<provider>.api_base`            | The API base URL for the provider.                                                                          | (Provider-specific)     |
| `<provider>.api_key`             | The API key for the provider.                                                                               |                          |
| `<provider>.model`               | The model name to use.                                                                                      | (Provider-specific)     |
| `<provider>.retries`             | The number of retry attempts for API requests.                                                              | `2`                     |
| `<provider>.proxy`               | The proxy URL to use (if needed).                                                                           |                          |
| `<provider>.max_tokens`          | The maximum number of tokens to generate.                                                                   | `2048`                   |
| `<provider>.top_p`               | The top-p value for nucleus sampling.                                                                       | `0.7`                    |
| `<provider>.temperature`         | The temperature value for controlling randomness.                                                            | `0.7`                    |
| `<provider>.frequency_penalty`   | The frequency penalty value.                                                                                | `0`                     |
| `<provider>.extra_headers`       | Extra headers to include in API requests (JSON string).                                                      | `{}`                    |
| `<provider>.completion_path`     | The API path for completion requests.                                                                      | (Provider-specific)     |
| `<provider>.answer_path`         | The JSON path to extract the answer from the API response.                                                   | (Provider-specific)     |
| `prompt.brief_commit_message`   | The prompt template for generating brief commit messages.                                                   | (See `defaults/defaults.go`) |
| `prompt.rich_commit_message`    | The prompt template for generating rich commit messages.                                                    | (See `defaults/defaults.go`) |
| `prompt.translation`             | The prompt template for translating commit messages.                                                         | (See `defaults/defaults.go`) |

**Note:** `<provider>` should be replaced with the actual provider name (e.g., `openai`, `gemini`, `claude`).

## Contribution

Contributions to GPTComet are welcome! Please refer to the [CONTRIBUTING.rst](CONTRIBUTING.rst) file for guidelines on how to contribute.

## Testing


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

-   This project is inspired by various AI-powered Git tools and utilizes several open-source libraries.
-   Thanks to the developers and maintainers of the LLMs used in this project.
-   The `gptcomet` directory contains a Python implementation of a similar tool, which served as a reference for certain features and concepts.
