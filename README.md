# Job Applier 3000

Centralised Web Application for finding employments for
KU Dept. of Computer Engineering Students and Alumni.

## How to install (and run)?

### Back-end

1. Install [golang](https://go.dev/)
1. `cd server`
1. `go get .`
1. `cp .env.example .env`
1. Replace the value with your MongoDB URI.
(alternatively, just move `mongo.go` somewhere else and skip this step).
1. `go run .`
1. Visit `localhost:8080/albums`. There should be json output.

Resources:
[Starter Code](https://go.dev/doc/tutorial/web-service-gin)
[Connecting to Database](https://www.slingacademy.com/article/securely-storing-secrets-with-environment-variables-in-go/)

### Front-end

1. Install [Bun](https://bun.sh)
1. `cd client`
1. `bun install` (according to ChatGPT).
1. `bun run dev --open`

Resources:
[Build an app with SvelteKit and Bun](https://bun.com/guides/ecosystem/sveltekit)
