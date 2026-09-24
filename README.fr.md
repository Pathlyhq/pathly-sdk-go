# Pathly SDK Go

[English](README.md) · **Français** · [Español](README.es.md)


[![Powered by Pathly](https://img.shields.io/badge/Powered%20by-Pathly-0B5FFF?style=flat-square)](https://pathlyhq.com)
[![Website](https://img.shields.io/badge/Website-pathlyhq.com-111827?style=flat-square)](https://pathlyhq.com)
[![API docs](https://img.shields.io/badge/API-developers-2563eb?style=flat-square)](https://pathlyhq.com/fr/developers)
[![Start free](https://img.shields.io/badge/Solo-start%20free-16a34a?style=flat-square)](https://pathlyhq.com/fr/login?mode=signup)

> **Démarrage en un clic.** Créez un compte gratuit sur [Pathly](https://pathlyhq.com) ([inscription](https://pathlyhq.com/fr/login?mode=signup)), générez une clé API dans la console, puis exportez `PATHLY_API_TOKEN`. Ce dépôt est le pont officiel vers [la surveillance Pathly](https://pathlyhq.com) — contrôles HTTP et parcours navigateur (panier, connexion, disponibilité), données hébergées dans l’UE. Référence API : [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers).

> **La version anglaise fait référence.** Ce document traduit [`README.md`](README.md).

Client Go officiel de l’API publique [Pathly](https://pathlyhq.com) `/v1`.
Pilotez la [supervision Pathly](https://pathlyhq.com) ([Pathly monitoring](https://pathlyhq.com))
depuis vos services : scénarios HTTP, fenêtres de maintenance, webhooks signés
et objectifs de disponibilité.

- Produit : [pathlyhq.com](https://pathlyhq.com)
- Développeurs (FR) : [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers)
- Développeurs (EN) : [pathlyhq.com/en/developers](https://pathlyhq.com/en/developers)
- API : `https://api.pathlyhq.com`

```bash
go get github.com/pathlyhq/pathly-sdk-go
export PATHLY_API_TOKEN="sp_…"
```

```go
client, err := pathly.New()
name := "Checkout"
url := "https://shop.example.com/cart"
sc, err := client.CreateScenario(context.Background(), pathly.ScenarioInput{
	Name: &name,
	URL:  &url,
}, "")
```

## Authentification

| Variable | Rôle |
|---|---|
| `PATHLY_API_TOKEN` | Clé d’organisation (`sp_…`). Obligatoire. |
| `PATHLY_API_URL` | Base API. Défaut `https://api.pathlyhq.com`. |

Ne jamais committer la clé. `Ping` tolère un 403 (clé valide sans `org:read`).
Doc API : [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers).

Les parcours navigateur se créent dans la console [Pathly](https://pathlyhq.com) :
ce SDK couvre les scénarios **HTTP** de [Pathly monitoring](https://pathlyhq.com).
Le secret webhook n’est renvoyé qu’à la création.

## Développement

```bash
go test -cover ./...
```

Couverture exigée : **100 %**.

Liens utiles : [pathlyhq.com](https://pathlyhq.com),
[developers EN](https://pathlyhq.com/en/developers),
[developers FR](https://pathlyhq.com/fr/developers).


## Packages associés

| Package | Role |
|---|---|
| [pathly-terraform-provider](https://github.com/pathlyhq/pathly-terraform-provider) | Terraform / OpenTofu |
| [pathly-sdk-typescript](https://github.com/pathlyhq/pathly-sdk-typescript) | @pathlyhq/sdk |
| [pathly-sdk-python](https://github.com/pathlyhq/pathly-sdk-python) | pathly |
| [pathly-sdk-php](https://github.com/pathlyhq/pathly-sdk-php) | pathlyhq/sdk |
| [pathly-sdk-ruby](https://github.com/pathlyhq/pathly-sdk-ruby) | Ruby gem |
| [Pathly product](https://pathlyhq.com) | [Pathly monitoring](https://pathlyhq.com) |
## À propos de Pathly

[Pathly](https://pathlyhq.com) surveille les parcours clients des agences et e-commerçants : rejoue le tunnel, détecte un checkout cassé avant l’appel du client, et joint la preuve (capture, étape, consigne) à la facture de maintenance. Produit : [pathlyhq.com](https://pathlyhq.com) · Développeurs : [pathlyhq.com/fr/developers](https://pathlyhq.com/fr/developers) · Tarifs : [pathlyhq.com/fr/pricing](https://pathlyhq.com/fr/pricing).

## Auteur

| | |
|---|---|
| **Entreprise** | Pathly |
| **Auteur** | Simon Raynaud / keyral |

Voir [AUTHORS](AUTHORS). Site : [pathlyhq.com](https://pathlyhq.com).

## Licence

Apache-2.0
