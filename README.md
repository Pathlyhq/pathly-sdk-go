# Pathly Go SDK

**English** · [Français](README.fr.md) · [Español](README.es.md)


[![Powered by Pathly](https://img.shields.io/badge/Powered%20by-Pathly-0B5FFF?style=flat-square)](https://pathlyhq.com)
[![Website](https://img.shields.io/badge/Website-pathlyhq.com-111827?style=flat-square)](https://pathlyhq.com)
[![API docs](https://img.shields.io/badge/API-developers-2563eb?style=flat-square)](https://pathlyhq.com/en/developers)
[![Start free](https://img.shields.io/badge/Solo-start%20free-16a34a?style=flat-square)](https://pathlyhq.com/en/login?mode=signup)

> **Get started in one click.** Create a free account on [Pathly](https://pathlyhq.com) ([sign up](https://pathlyhq.com/en/login?mode=signup)), create an API key in the console, then export `PATHLY_API_TOKEN`. This project is the official bridge to [Pathly monitoring](https://pathlyhq.com) — real-browser and HTTP checks for checkout, login and availability, with data hosted in the EU. Full API reference: [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers).

Official Go client for the [Pathly](https://pathlyhq.com) public `/v1` API.
Wire [Pathly monitoring](https://pathlyhq.com) into your services, workers and
CI: HTTP scenarios, maintenance windows, signed webhooks and SLA targets.

- Product: [pathlyhq.com](https://pathlyhq.com)
- Developers (EN): [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers)
- Developers (FR): [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers)
- API host: `https://api.pathlyhq.com`

```bash
go get github.com/pathlyhq/pathly-sdk-go
export PATHLY_API_TOKEN="sp_…"
```

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pathlyhq/pathly-sdk-go"
)

func main() {
	client, err := pathly.New() // reads PATHLY_API_TOKEN
	if err != nil {
		log.Fatal(err)
	}
	name := "Checkout"
	url := "https://shop.example.com/cart"
	interval := int64(300)
	sc, err := client.CreateScenario(context.Background(), pathly.ScenarioInput{
		Name:        &name,
		URL:         &url,
		IntervalSec: &interval,
	}, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sc.ID)
}
```

## Why Pathly monitoring

[Pathly](https://pathlyhq.com) is EU-hosted synthetic monitoring for the journeys
your customers actually take. Use this SDK when you want the same [Pathly
monitoring](https://pathlyhq.com) controls as the console, from Go:

| Need | Link |
|---|---|
| Product overview | [pathlyhq.com](https://pathlyhq.com) |
| API & SDK docs (EN) | [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers) |
| API & SDK docs (FR) | [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers) |
| Status & uptime story | [Pathly monitoring](https://pathlyhq.com) on [pathlyhq.com](https://pathlyhq.com) |

## Authentication

| Variable | Purpose |
|---|---|
| `PATHLY_API_TOKEN` | Organization API key (`sp_` prefix). Required. |
| `PATHLY_API_URL` | API base. Defaults to `https://api.pathlyhq.com`. |

Never hard-code the token. Prefer the environment or a secret store.

`Client.Ping` calls `GET /v1/usage` and **accepts HTTP 403**: the key is valid
but lacks `org:read`. That keeps least privilege for scenario-only keys.

Scopes: `scenarios:read/write`, `alerting:read/write`, `maintenance:read/write`,
`sla:read/write`. Details on [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers).

## Resources

| Resource | Methods |
|---|---|
| Scenario (HTTP) | `CreateScenario`, `GetScenario`, `UpdateScenario`, `DeleteScenario`, `ListScenarios`, `MuteScenario` |
| Maintenance window | `CreateMaintenanceWindow`, `GetMaintenanceWindow`, `DeleteMaintenanceWindow`, `ListMaintenanceWindows` |
| Webhook | `CreateWebhook`, `GetWebhook`, `DeleteWebhook`, `ListWebhooks` |
| SLA target | `UpsertSlaTarget` (PUT), `GetSlaTarget`, `DeleteSlaTarget`, `ListSlaTargets` |

Browser journeys are **not** managed here: create them in the [Pathly](https://pathlyhq.com)
console. This package is for **HTTP** [Pathly monitoring](https://pathlyhq.com)
scenarios only.

The webhook **secret** is returned once at creation. Store it in a vault. The
API never returns the destination URL on read, only `urlFingerprint`.

## Behavior

- Creates send an `Idempotency-Key` (auto-generated unless you pass one).
- HTTP 429 and 5xx honour `Retry-After` (capped at 90 seconds, four attempts).
- HTTP 404 is detectable with `pathly.IsNotFound(err)`.
- Zero third-party runtime dependencies (stdlib only).

## Development

```bash
go test -cover ./...
```

Coverage is required at **100%**.

See [examples/complete](./examples/complete) and the API reference on
[pathlyhq.com/en/developers](https://pathlyhq.com/en/developers)
([FR](https://pathlyhq.com/fr/developers)).

## Related packages

| Package | Role |
|---|---|
| [pathly-terraform-provider](https://github.com/pathlyhq/pathly-terraform-provider) | Terraform / OpenTofu |
| [pathly-sdk-typescript](https://github.com/pathlyhq/pathly-sdk-typescript) | `@pathlyhq/sdk` |
| [pathly-sdk-python](https://github.com/pathlyhq/pathly-sdk-python) | `pathly` |
| [pathly-sdk-php](https://github.com/pathlyhq/pathly-sdk-php) | `pathlyhq/sdk` |
| [Pathly product](https://pathlyhq.com) | [Pathly monitoring](https://pathlyhq.com) |

## About Pathly

[Pathly](https://pathlyhq.com) is synthetic monitoring for agencies and e-commerce: replay the customer journey, catch broken checkouts before your clients call, and keep evidence (screenshot, step, runbook) ready for the invoice. Product: [pathlyhq.com](https://pathlyhq.com) · Developers: [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers) · Status & pricing: [pathlyhq.com/en/pricing](https://pathlyhq.com/en/pricing).

## Author

| | |
|---|---|
| **Company** | Pathly |
| **Author** | Simon Raynaud / keyral |

See [AUTHORS](AUTHORS). Homepage: [pathlyhq.com](https://pathlyhq.com).

## License

Apache-2.0
