# likho

Likho is a lightweight, fast, and simple static site generator written in Go. It's designed to convert markdown files into a fully functional static website with minimal configuration.

## Features

- Markdown to HTML conversion
- YAML-based configuration
- Custom metadata for each post
- Automatic sitemap and RSS feed generation
- Command-line interface for easy management

## Installation

1. Ensure you have Go 1.23 or later installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/intothevoid/likho.git
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
   ./likho create "My New Post" [flags]
   ```
   
   Available flags:
   - `-t, --tags string`: Comma-separated list of tags for the post
   - `-i, --image string`: URL of the featured image for the post

   Example:
   ```
   ./likho create "My New Post" -t "technology,golang" -i "https://example.com/image.jpg"
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

```yaml
site:
  title: "My Awesome Blog"
  description: "A blog about awesome things"
  baseURL: "https://example.com"

content:
  sourceDir: "content"
  postsDir: "posts"

output:
  publicDir: "public"

templates:
  dir: "templates"

static:
  dir: "static"

```
### Configuration options:

- `site`:
  - `title`: The title of your blog
  - `description`: A brief description of your blog
  - `baseURL`: The base URL where your site will be hosted

- `content`:
  - `sourceDir`: The directory containing your content files
  - `postsDir`: The subdirectory within `sourceDir` where posts are stored

- `output`:
  - `publicDir`: The directory where generated HTML files will be placed

- `templates`:
  - `dir`: The directory containing your HTML templates

- `static`:
  - `dir`: The directory containing static assets (CSS, images, etc.)

Ensure that your `config.yaml` file is in the root directory of your Likho project.