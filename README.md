# Likho

<a href="assets/logo.jpg"><img src="assets/logo.jpg" height="300px" alt="Likho Logo"></a>

[![Go](https://github.com/intothevoid/likho/actions/workflows/go.yml/badge.svg)](https://github.com/intothevoid/likho/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/intothevoid/likho)](https://goreportcard.com/report/github.com/intothevoid/likho)
[![Release](https://img.shields.io/github/release/intothevoid/likho.svg)](https://github.com/intothevoid/likho/releases/latest)
[![License](https://img.shields.io/github/license/intothevoid/likho.svg)](https://github.com/intothevoid/likho/blob/main/LICENSE)

Likho is a lightweight, fast, and simple static site generator written in Go. 

It's designed to convert markdown files into a fully functional static website with minimal configuration.

# Demo

A demo site (my personal blog) built with Likho can be found at [https://intothevoid.github.io/](https://intothevoid.github.io/)

## Features

- Markdown to HTML conversion
- YAML-based configuration
- Custom metadata for each post and page
- Automatic sitemap and RSS feed generation
- Command-line interface for easy management
- Support for creating both posts and pages

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

### Create a new post

```
./likho create post "My New Post" [flags]
```

Available flags:
- `-t, --tags string`: Comma-separated list of tags for the post
- `-i, --image string`: URL of the featured image for the post
- `-d, --description string`: Short description of the post

Example:
```
./likho create post "My New Post" -t "technology,golang" -i "https://example.com/image.jpg"
```

### Create a new page

```
./likho create page "Page Title" [flags]
```

Available flags:
- `-i, --image string`: URL of the featured image for the page
- `-d, --description string`: Short description of the page

Example:
```
./likho create page "About Me" -i "https://example.com/about.jpg" -d "Learn more about the author"
```

### Generate the static site

```
./likho generate
```

This command will:
1. Parse all posts and pages
2. Generate HTML files for each post and page
3. Create an index page with recent posts
4. Generate a sitemap and RSS feed
5. Copy the CSS file to the output directory

### Serve the generated site locally

```
./likho serve
```

This command starts a local web server to preview your generated site.

### Display help information

```
./likho help
```

Use this command to see all available commands and their descriptions.

## Generate with Docker

Use the following command to build a Docker image:

```
docker build -t likho .
```

To run the Docker container and generate the static site, use:

```
docker run --rm -v "$(pwd)/content:/app/content:ro" -v "$(pwd)/public:/app/public" likho:latest generate
```

Note: The following directories are required:
- content - Contains the markdown files for the site. This is the input directory.
- public - Will contain the generated static site. This is the output directory.

### Push to Docker Hub

```
docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t yourdockerhubusername/likho:latest --push .
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

- `assets`:
  - `dir`: The directory containing static assets (CSS, images, etc.)

Ensure that your `config.yaml` file is in the root directory of your Likho project.

## Directory Structure

After setting up your Likho project, your directory structure should look like this:

```
my-likho-site/
├── config.yaml
├── assets/
│   └── logo.jpg
|   └── main.css
├── content/
│   ├── posts/
│   │   └── YYYY-MM-DD/
│   │       └── post-slug.md
│   ├── pages/
│   │   └── page-slug.md
├── templates/
│   ├── base.html
│   ├── index.html
│   ├── post.html
│   ├── pages.html
│   ├── header.html
│   └── footer.html
└── public/
    └── (generated files)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.