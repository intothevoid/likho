# likho

Likho is a lightweight, fast, and simple static site generator written in Go. It's designed to convert markdown files into a fully functional static website with minimal configuration.

## Features

- Markdown to HTML conversion
- YAML-based configuration
- Custom metadata for each post
- Automatic sitemap and RSS feed generation
- Command-line interface for easy management

## Installation

1. Ensure you have Go 1.22 or later installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/yourusername/likho.git
   cd likho
   ```

3. Build the application:
   ```
   go build -o likho cmd/likho/main.go
   ```

4. (Optional) Add the `likho` binary to your PATH for easier access.

## Usage

Likho provides several commands to manage your static site:

1. Create a new post:
   ```
   ./likho create "My New Post"
   ```

2. Generate the static site:
   ```
   ./likho generate
   ```

3. Serve the generated site locally:
   ```
   ./likho serve
   ```

4. Display help information:
   ```
   ./likho help
   ```

## Configuration

Likho uses a `config.yaml` file in the root directory for site-wide configuration. Here's an example:
