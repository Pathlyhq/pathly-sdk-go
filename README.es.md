# Pathly Go SDK

[English](README.md) · [Français](README.fr.md) · **Español**


[![Powered by Pathly](https://img.shields.io/badge/Powered%20by-Pathly-0B5FFF?style=flat-square)](https://pathlyhq.com)
[![Website](https://img.shields.io/badge/Website-pathlyhq.com-111827?style=flat-square)](https://pathlyhq.com)
[![API docs](https://img.shields.io/badge/API-developers-2563eb?style=flat-square)](https://pathlyhq.com/es/developers)
[![Start free](https://img.shields.io/badge/Solo-start%20free-16a34a?style=flat-square)](https://pathlyhq.com/es/login?mode=signup)

> **Empiece en un clic.** Cree una cuenta gratuita en [Pathly](https://pathlyhq.com) ([registro](https://pathlyhq.com/es/login?mode=signup)), genere una clave API en la consola y exporte `PATHLY_API_TOKEN`. Este repositorio es el puente oficial hacia [la monitorización Pathly](https://pathlyhq.com): comprobaciones HTTP y de navegador (carrito, login, disponibilidad), con datos en la UE. Referencia API: [pathlyhq.com/es/developers](https://pathlyhq.com/es/developers).

> **La versión en inglés es la referencia.** Este documento traduce [`README.md`](README.md).

Cliente Go oficial de la API pública [Pathly](https://pathlyhq.com) `/v1`.
Integre [Pathly monitoring](https://pathlyhq.com) en sus servicios, workers y
CI: escenarios HTTP, ventanas de mantenimiento, webhooks firmados y objetivos SLA.

- Producto: [pathlyhq.com](https://pathlyhq.com)
- Desarrolladores (ES): [pathlyhq.com/es/developers](https://pathlyhq.com/es/developers)
- Desarrolladores (EN): [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers)
- Desarrolladores (FR): [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers)
- Host API: `https://api.pathlyhq.com`

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

## Por qué Pathly monitoring

[Pathly](https://pathlyhq.com) es monitorización sintética alojada en la UE para
los recorridos que realmente hacen sus clientes. Use este SDK cuando quiera los
mismos controles de [Pathly monitoring](https://pathlyhq.com) que en la consola,
desde Go:

| Necesidad | Enlace |
|---|---|
| Visión del producto | [pathlyhq.com](https://pathlyhq.com) |
| Docs API y SDK (ES) | [pathlyhq.com/es/developers](https://pathlyhq.com/es/developers) |
| Docs API y SDK (EN) | [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers) |
| Docs API y SDK (FR) | [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers) |
| Estado y disponibilidad | [Pathly monitoring](https://pathlyhq.com) en [pathlyhq.com](https://pathlyhq.com) |

## Autenticación

| Variable | Propósito |
|---|---|
| `PATHLY_API_TOKEN` | Clave de API de la organización (prefijo `sp_`). Obligatoria. |
| `PATHLY_API_URL` | Base de la API. Por defecto `https://api.pathlyhq.com`. |

Nunca codifique el token. Prefiera el entorno o un almacén de secretos.

`Client.Ping` llama a `GET /v1/usage` y **acepta HTTP 403**: la clave es válida
pero carece de `org:read`. Así se mantiene el menor privilegio para claves solo
de escenarios.

Scopes: `scenarios:read/write`, `alerting:read/write`, `maintenance:read/write`,
`sla:read/write`. Detalles en [pathlyhq.com/es/developers](https://pathlyhq.com/es/developers).

## Recursos

| Recurso | Métodos |
|---|---|
| Escenario (HTTP) | `CreateScenario`, `GetScenario`, `UpdateScenario`, `DeleteScenario`, `ListScenarios`, `MuteScenario` |
| Ventana de mantenimiento | `CreateMaintenanceWindow`, `GetMaintenanceWindow`, `DeleteMaintenanceWindow`, `ListMaintenanceWindows` |
| Webhook | `CreateWebhook`, `GetWebhook`, `DeleteWebhook`, `ListWebhooks` |
| Objetivo SLA | `UpsertSlaTarget` (PUT), `GetSlaTarget`, `DeleteSlaTarget`, `ListSlaTargets` |

Los recorridos de navegador **no** se gestionan aquí: créelos en la consola
[Pathly](https://pathlyhq.com). Este paquete es solo para escenarios **HTTP**
de [Pathly monitoring](https://pathlyhq.com).

El **secret** del webhook se devuelve una sola vez en la creación. Guárdelo en
un vault. La API nunca devuelve la URL de destino en lectura, solo
`urlFingerprint`.

## Comportamiento

- Las creaciones envían un `Idempotency-Key` (generado automáticamente salvo que pase uno).
- HTTP 429 y 5xx respetan `Retry-After` (tope de 90 segundos, cuatro intentos).
- HTTP 404 se detecta con `pathly.IsNotFound(err)`.
- Sin dependencias de runtime de terceros (solo stdlib).

## Desarrollo

```bash
go test -cover ./...
```

Se exige cobertura al **100%**.

Véase [examples/complete](./examples/complete) y la referencia API en
[pathlyhq.com/es/developers](https://pathlyhq.com/es/developers)
([EN](https://pathlyhq.com/en/developers), [FR](https://pathlyhq.com/fr/developers)).

## Paquetes relacionados

| Paquete | Rol |
|---|---|
| [pathly-terraform-provider](https://github.com/pathlyhq/pathly-terraform-provider) | Terraform / OpenTofu |
| [pathly-sdk-typescript](https://github.com/pathlyhq/pathly-sdk-typescript) | `@pathlyhq/sdk` |
| [pathly-sdk-python](https://github.com/pathlyhq/pathly-sdk-python) | `pathly` |
| [pathly-sdk-php](https://github.com/pathlyhq/pathly-sdk-php) | `pathlyhq/sdk` |
| [Pathly product](https://pathlyhq.com) | [Pathly monitoring](https://pathlyhq.com) |

## Acerca de Pathly

[Pathly](https://pathlyhq.com) es monitorización sintética para agencias y e-commerce: reproduce el recorrido del cliente, detecta un checkout roto antes de la llamada, y deja la prueba (captura, paso, runbook) lista para la factura. Producto: [pathlyhq.com](https://pathlyhq.com) · Desarrolladores: [pathlyhq.com/es/developers](https://pathlyhq.com/es/developers) · Precios: [pathlyhq.com/es/pricing](https://pathlyhq.com/es/pricing).

## Autor

| | |
|---|---|
| **Empresa** | Pathly |
| **Autor** | Simon Raynaud / keyral |

Véase [AUTHORS](AUTHORS). Sitio: [pathlyhq.com](https://pathlyhq.com).

## Licencia

Apache-2.0
