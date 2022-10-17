# Weblog

![img](https://cln.sh/wugjnghoVh5ytuqq8xXn/download)

Weblog is a minimal blog that will serve as a real-world example project for Bud.

Blogs are expecially good examples because they're fairly easy to build, but demonstrate a lot of the fundamental capabilities you'd need to build any website.

Much of the code in here will make its way into Bud either through runtime libraries or code generation.

This initial blog took me about 8 hours to write. I hope in the future with Bud, the same blog will take an hour or two. I plan to slowly transition this repo over to Bud as required features are added.

Since this repository is acting as a "canary" for future features, **[please share your feedback](https://github.com/livebud/weblog/issues/new)** on areas in the codebase that you find confusing or don't like! Also if you find features in here that aren't in Bud, consider contributing them to Bud with a PR!

## Features

- Controllers
- Models
- Views
- Templates (html/template)
- Custom Routing
- Migrations
- Database access
- Middleware
- Authentication
- CSRF protection

## Install

**Prerequisite**

1. [`direnv`](https://direnv.net/docs/installation.html)
2. go >= 1.18
3. Install the following tools to your go path

```sh
$ go install github.com/matthewmueller/migrate/cmd/migrate
$ go install github.com/matthewmueller/pogo/cmd/pogo
```

**Setup**

```sh
# Install go dependencies
go mod tidy

# Create a postgres database (assumes you have a Postgres database running)
createdb weblog

# Allow the .envrc to add variables to your project environment
direnv allow

# Migrate the database
make migrate.up

# Run the weblog
go run bud/cmd/app/main.go
```

## License

MIT
