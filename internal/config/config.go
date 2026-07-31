package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the entire application configuration, serialized to YAML.
type Config struct {
	Title        string     `yaml:"title"`
	DefaultView  string     `yaml:"default_view"`
	DefaultTheme string     `yaml:"default_theme"`
	PasswordHash string     `yaml:"password_hash,omitempty"`
	Health       Health     `yaml:"health"`
	Appearance   Appearance `yaml:"appearance"`
	Sections     []Section  `yaml:"sections"`
	Links        []Link     `yaml:"links"`
}

type Health struct {
	Enabled         bool `yaml:"enabled"`
	IntervalSeconds int  `yaml:"interval_seconds"`
	TimeoutSeconds  int  `yaml:"timeout_seconds"`
}

type Appearance struct {
	Page      Page      `yaml:"page"`
	Grid      Grid      `yaml:"grid"`
	Item      Item      `yaml:"item"`
	Icon      IconSizes `yaml:"icon"`
	Text      Text      `yaml:"text"`
	StatusDot StatusDot `yaml:"status_dot"`
}

type Page struct {
	MaxWidth   int    `yaml:"max_width"`
	Background string `yaml:"background"`
	FontFamily string `yaml:"font_family"`
	FontSize   int    `yaml:"font_size"`
}

type Grid struct {
	Columns      string `yaml:"columns"`
	MinItemWidth int    `yaml:"min_item_width"`
	Gap          int    `yaml:"gap"`
}

type Item struct {
	CornerRadius   int  `yaml:"corner_radius"`
	Padding        int  `yaml:"padding"`
	Background     bool `yaml:"background"`
	Border         bool `yaml:"border"`
	BorderStrength int  `yaml:"border_strength"`
	Shadow         bool `yaml:"shadow"`
	ShadowStrength int  `yaml:"shadow_strength"`
}

type IconSizes struct {
	SizeDefault int `yaml:"size_default"`
	SizeCompact int `yaml:"size_compact"`
	SizeCard    int `yaml:"size_card"`
	SizeLarge   int `yaml:"size_large"`
}

type Text struct {
	Align           string `yaml:"align"`
	ShowDescription bool   `yaml:"show_description"`
}

type StatusDot struct {
	Enabled  bool   `yaml:"enabled"`
	Size     int    `yaml:"size"`
	Position string `yaml:"position"`
}

type Section struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Link struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	URL         string `yaml:"url"`
	Icon        string `yaml:"icon"`
	Color       string `yaml:"color"`
	Section     string `yaml:"section,omitempty"`
	Health      *bool  `yaml:"health,omitempty"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	t := true
	return Config{
		Title:        "Minidash",
		DefaultView:  "default",
		DefaultTheme: "auto",
		Health:       Health{Enabled: true, IntervalSeconds: 60, TimeoutSeconds: 5},
		Appearance: Appearance{
			Page:      Page{MaxWidth: 1200, Background: "", FontFamily: "system", FontSize: 16},
			Grid:      Grid{Columns: "auto", MinItemWidth: 220, Gap: 16},
			Item:      Item{CornerRadius: 12, Padding: 16, Background: true, Border: true, BorderStrength: 1, Shadow: true, ShadowStrength: 1},
			Icon:      IconSizes{SizeDefault: 24, SizeCompact: 20, SizeCard: 32, SizeLarge: 96},
			Text:      Text{Align: "left", ShowDescription: true},
			StatusDot: StatusDot{Enabled: true, Size: 8, Position: "bottom-right"},
		},
		Links: []Link{{Name: "Minidash", URL: "https://example.com", Icon: "lucide:home", Color: "#4f9cff", Health: &t}},
	}
}

// Load reads config from path. Missing file -> writes & returns Default().
// Malformed file -> error (caller keeps last-good in-memory copy).
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c := Default()
			if e := Save(path, &c); e != nil {
				return nil, e
			}
			return &c, nil
		}
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	c.applyDefaults()
	return &c, nil
}

func (c *Config) applyDefaults() {
	d := Default()
	if c.Title == "" {
		c.Title = d.Title
	}
	if c.DefaultView == "" {
		c.DefaultView = d.DefaultView
	}
	if c.DefaultTheme == "" {
		c.DefaultTheme = d.DefaultTheme
	}
	if c.Health.IntervalSeconds == 0 {
		c.Health.IntervalSeconds = d.Health.IntervalSeconds
	}
	if c.Health.TimeoutSeconds == 0 {
		c.Health.TimeoutSeconds = d.Health.TimeoutSeconds
	}
	if c.Appearance.Page.MaxWidth == 0 {
		c.Appearance.Page = d.Appearance.Page
	}
	if c.Appearance.Grid.MinItemWidth == 0 && c.Appearance.Grid.Gap == 0 {
		c.Appearance.Grid = d.Appearance.Grid
	}
	if c.Appearance.Item.CornerRadius == 0 && c.Appearance.Item.Padding == 0 {
		c.Appearance.Item = d.Appearance.Item
	}
	if c.Appearance.Icon.SizeLarge == 0 {
		c.Appearance.Icon = d.Appearance.Icon
	}
	if c.Appearance.Text.Align == "" {
		c.Appearance.Text = d.Appearance.Text
	}
	if c.Appearance.StatusDot.Size == 0 {
		c.Appearance.StatusDot = d.Appearance.StatusDot
	}
}

// Save writes config atomically: back up current file, write temp, rename.
func Save(path string, c *Config) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if prev, e := os.ReadFile(path); e == nil {
		_ = os.WriteFile(path+".bak", prev, 0o644)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config.yaml.*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
