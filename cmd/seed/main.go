package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"minidash/internal/config"
	"minidash/internal/icons"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	path := os.Getenv("MINIDASH_CONFIG")
	if path == "" {
		path = "config.yaml"
	}

	cfg, err := config.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load:", err)
		os.Exit(1)
	}

	// Build a pool of valid icons (rendered names) from the icon packs.
	iconPool := []string{}
	for _, pfx := range []string{"lucide", "simple-icons", "tabler"} {
		for _, r := range icons.Search("", pfx, 300) {
			iconPool = append(iconPool, r.Prefix+":"+r.Name)
		}
	}
	if len(iconPool) == 0 {
		fmt.Fprintln(os.Stderr, "no icons found")
		os.Exit(1)
	}
	randIcon := func() string { return iconPool[rng.Intn(len(iconPool))] }

	palette := []string{"#4f9cff", "#62a9ff", "#a6d189", "#fabd2f", "#f38ba8", "#bd93f9",
		"#8caaee", "#94b386", "#ff7755", "#22c55e", "#ef4444", "#e57373", "#88c0d0",
		"#50fa7b", "#ff79c6", "#f1fa8c", "#8be9fd", "#ffb86c", "#6272a4", "#fab387"}
	randColor := func() string {
		if rng.Intn(2) == 0 {
			return palette[rng.Intn(len(palette))]
		}
		return fmt.Sprintf("#%02x%02x%02x", 90+rng.Intn(150), 90+rng.Intn(150), 90+rng.Intn(150))
	}

	descs := []string{
		"A handy tool for your homelab.", "Self-hosted and open source.", "Runs in Docker.",
		"Sits behind the reverse proxy.", "Organize your digital life.", "Fast, lightweight, reliable.",
		"Web UI included.", "Sync across all your devices.", "Monitor and automate everything.",
		"Community-driven project.", "Scales nicely.", "Set-and-forget daemon.", "Your personal gateway.",
		"", "", "", // ~1/4 end up with no description (omitempty)
	}
	randDesc := func() string { return descs[rng.Intn(len(descs))] }

	domains := []string{"example.lan", "homelab.local", "lab.net", "internal", "dev"}
	slugRe := regexp.MustCompile(`[^a-z0-9]+`)
	slug := func(s string) string {
		return strings.Trim(slugRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
	}
	randURL := func(name string) string {
		scheme := "https"
		if rng.Intn(3) == 0 {
			scheme = "http"
		}
		return fmt.Sprintf("%s://%s.%s", scheme, slug(name), domains[rng.Intn(len(domains))])
	}

	off := false
	noHealth := func() *bool { return &off } // seeded items are not health-monitored

	// Sections to add (skip IDs that already exist).
	type secDef struct{ id, name string }
	sectionDefs := []secDef{
		{"media", "Media"}, {"homelab", "Home Lab"}, {"dev", "Dev Tools"}, {"social", "Social"},
		{"docs", "Docs & Wiki"}, {"monitoring", "Monitoring"}, {"downloads", "Downloads"},
		{"network", "Network"}, {"ai", "AI & ML"}, {"games", "Gaming"},
	}
	existing := map[string]bool{}
	for _, s := range cfg.Sections {
		existing[s.ID] = true
	}
	sectionByID := map[string]string{}
	for _, s := range sectionDefs {
		sectionByID[s.id] = s.name
		if !existing[s.id] {
			cfg.Sections = append(cfg.Sections, config.Section{ID: s.id, Name: s.name})
		}
	}
	allSecIDs := []string{}
	for _, s := range sectionDefs {
		allSecIDs = append(allSecIDs, s.id)
	}
	randSection := func() string { return allSecIDs[rng.Intn(len(allSecIDs))] }

	// Real-ish service names per category.
	services := map[string][]string{
		"media":      {"Plex", "Jellyfin", "Emby", "Sonarr", "Radarr", "Bazarr", "Prowlarr", "Tautulli", "Overseerr", "Navidrome", "Komga", "Kavita", "TubeArchivist"},
		"homelab":    {"Home Assistant", "Portainer", "Yacht", "Nginx Proxy Manager", "Traefik", "Caddy", "Proxmox", "CasaOS", "Umbrel", "Homarr", "Dashy", "Cockpit"},
		"dev":        {"Gitea", "Forgejo", "Drone CI", "Woodpecker CI", "Jenkins", "Vault", "Docker Registry", "Verdaccio", "SonarQube", "Code-Server", "HedgeDoc"},
		"social":     {"Mastodon", "Lemmy", "Pixelfed", "Discord", "Telegram", "Matrix", "Synapse", "Rocket Chat", "Friendica"},
		"docs":       {"Wiki.js", "Outline", "BookStack", "DokuWiki", "TiddlyWiki", "Obsidian", "AppFlowy", "Affine", "Memos"},
		"monitoring": {"Grafana", "Prometheus", "Loki", "Alertmanager", "Uptime Kuma", "Netdata", "Beszel", "Zabbix", "Speedtest Tracker"},
		"downloads":  {"qBittorrent", "Transmission", "Deluge", "SABnzbd", "NZBGet", "Jackett", "Metube", "Aria2", "Prowlarr"},
		"network":    {"Pi-hole", "AdGuard Home", "Unbound", "WireGuard", "Tailscale", "OpenVPN", "Cloudflared", "Nginx", "pfSense", "OPNsense", "UniFi"},
		"ai":         {"Ollama", "Open WebUI", "ComfyUI", "Stable Diffusion", "LocalAI", "AnythingLLM", "Dify", "Flowise", "n8n"},
		"games":      {"Minecraft", "Valheim", "Terraria", "Factorio", "Palworld", "CSGO", "ARK", "Rust", "Pterodactyl"},
	}

	addLink := func(name, section string) {
		cfg.Links = append(cfg.Links, config.Link{
			Type:        "link",
			Name:        name,
			Description: randDesc(),
			URL:         randURL(name),
			Icon:        randIcon(),
			Color:       randColor(),
			Section:     section,
			Health:      noHealth(),
		})
	}

	// 1) One link per service, in its category section.
	for cat, names := range services {
		for _, n := range names {
			addLink(n, cat)
		}
	}

	// 2) Lots of filler links: spread across all sections + a chunk unsectioned.
	adjectives := []string{"Alpha", "Beta", "Quantum", "Cosmic", "Neon", "Cyber", "Solar", "Lunar",
		"Aurora", "Crimson", "Azure", "Emerald", "Onyx", "Cobalt", "Frost", "Storm", "Shadow", "Nova"}
	nouns := []string{"Node", "Cluster", "Vault", "Portal", "Hub", "Engine", "Matrix", "Relay",
		"Beacon", "Gateway", "Forge", "Nexus", "Citadel", "Bastion", "Oracle", "Sentinel", "Proxy", "Bridge"}
	for i := 0; i < 90; i++ {
		name := adjectives[rng.Intn(len(adjectives))] + " " + nouns[rng.Intn(len(nouns))]
		sec := ""
		if rng.Intn(4) != 0 {
			sec = randSection() // ~75% sectioned
		}
		addLink(name, sec)
	}

	// 3) Notes spread across sections (some with a title, some without).
	noteTitles := []string{"", "", "Reminder", "Tip", "Note", "TODO", "Credentials", "How-to", ""}
	noteBodies := []string{
		"Back up the NAS before the weekend.", "SSH key is in the vault.", "Restart the proxy after cert renewal.",
		"DNS: 10.0.0.1 primary, 9.9.9.9 fallback.", "qBittorrent webui on :8080.", "WireGuard port 51820.",
		"Daily backup runs at 03:00.", "Don't forget to renew the domain.", "Grafana admin password in Bitwarden.",
		"ComfyUI model path: /models/stable-diffusion", "Plex transcode dir is on the SSD.", "Update containers monthly.",
		"Pihole gravity update every sunday.", "Test restore quarterly.", "Labels go on every cable.",
	}
	for i := 0; i < 24; i++ {
		body := noteBodies[rng.Intn(len(noteBodies))]
		cfg.Links = append(cfg.Links, config.Link{
			Type:    "note",
			Name:    noteTitles[rng.Intn(len(noteTitles))],
			Text:    body,
			URL:     "",
			Icon:    randIcon(),
			Color:   randColor(),
			Section: randSection(),
		})
	}

	if err := config.Save(path, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "save:", err)
		os.Exit(1)
	}
	fmt.Printf("Seeded %d sections, %d links/notes total into %s\n", len(cfg.Sections), len(cfg.Links), path)
}
