# Job Applier 3000

[![Run Vitest with Bun](https://github.com/lnwdevelopers007/job-applier-3000/actions/workflows/vitest.yaml/badge.svg)](https://github.com/lnwdevelopers007/job-applier-3000/actions/workflows/vitest.yaml)
[![Go Tests](https://github.com/lnwdevelopers007/job-applier-3000/actions/workflows/go.yaml/badge.svg)](https://github.com/lnwdevelopers007/job-applier-3000/actions/workflows/go.yaml)

Centralised Web Application for finding employments for
KU Dept. of Computer Engineering Students and Alumni.

## How to install (and run)?

### Quick run (No language/compiler installation needed)

1. Setup configurations for .env from .env.example provided in 3 places (Outer most, Backend and Frontend)
1. For Backend and Frontend, copy your .env to .env.production and read .env.example as standard.
1. Run `docker-compose up --build`
1. Run `docker-compose down` to stop the services.

### Back-end

1. Install [golang](https://go.dev/)
1. `cd server`
1. `go get .`
1. `cp .env.example .env`
1. Replace it with your real config.
1. `go run .`
1. Visit `localhost:8080/jobs`.
1. To run tests, use `godotenv -f ./.env go test ./... -v`
   (install [godotenv](https://github.com/joho/godotenv?tab=readme-ov-file#installation) as bin command first.)

Resources:

-   [Starter Code](https://go.dev/doc/tutorial/web-service-gin)
-   [Connecting to Database](https://www.slingacademy.com/article/securely-storing-secrets-with-environment-variables-in-go/)

### Front-end

1. Install [Bun](https://bun.sh)
1. `cd client`
1. `cp .env.example .env`
1. Replace it with your real config.
1. `bun install`
1. `bun run dev --open`

Resources:
[Build an app with SvelteKit and Bun](https://bun.com/guides/ecosystem/sveltekit)
