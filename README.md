# hexclock

A tiny Go HTTP service that reports the current UTC time as JSON.

| Endpoint       | Response                                                              |
| -------------- | --------------------------------------------------------------------- |
| `GET /`        | `{"utc":"2026-10-03T17:34:56Z","unix":1791048896,"hex":"0x6ac13cc0"}` |
| `GET /healthz` | `{"status":"ok"}`                                                     |

The server listens on port `8080`. Set `PORT` to change it.

## Run the published image

```sh
docker run --rm -p 8080:8080 ghcr.io/fregas/agentbench-hexclock-v3-17-b:latest
curl localhost:8080/
curl localhost:8080/healthz
```

Pin a release with a version tag, for example `:0.1.1`.

## Develop

```sh
go test ./...                    # unit tests
docker build -t hexclock .       # multi-stage build; runs vet + tests, ships a distroless non-root image
```

## Releases

Pushing a tag that matches `v*` (for example `v0.2.0`) runs
`.github/workflows/release-image.yml`. The workflow builds the image and pushes
it to `ghcr.io/fregas/agentbench-hexclock-v3-17-b` with the tags `0.2.0` and
`latest`. It authenticates with the workflow's own `GITHUB_TOKEN`.
