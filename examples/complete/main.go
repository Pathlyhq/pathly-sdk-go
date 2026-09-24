package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pathlyhq/pathly-sdk-go"
)

// Complete walkthrough against the live Pathly API.
// Requires PATHLY_API_TOKEN. See https://pathlyhq.com/en/developers
func main() {
	client, err := pathly.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	name := "sdk-go-example"
	u := "https://example.com/"
	interval := int64(3600)
	sc, err := client.CreateScenario(ctx, pathly.ScenarioInput{
		Name:        &name,
		URL:         &u,
		IntervalSec: &interval,
	}, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("scenario", sc.ID)

	wh, err := client.CreateWebhook(ctx, pathly.WebhookInput{
		URL:    "https://hooks.example.com/pathly",
		Events: []string{"run.failed", "run.recovered"},
	}, "")
	if err != nil {
		log.Fatal(err)
	}
	if wh.Secret != nil {
		fmt.Println("webhook secret (store once):", *wh.Secret)
	}

	mon := sc.ID
	slaName := "example 99.9"
	_, err = client.UpsertSlaTarget(ctx, pathly.SlaTargetInput{
		MonitorID:    &mon,
		Name:         &slaName,
		ObjectivePct: 99.9,
		WindowDays:   30,
	}, "")
	if err != nil {
		log.Fatal(err)
	}

	_ = client.DeleteWebhook(ctx, wh.ID)
	_ = client.DeleteScenario(ctx, sc.ID)
	fmt.Println("cleaned up — Pathly monitoring via", os.Getenv("PATHLY_API_URL"))
}
