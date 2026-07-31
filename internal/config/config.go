package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the entire application configuration, serialized to YAML.
type Config struct {
	Title        string     `yaml:"title" json:"title"`
	DefaultView  string     `yaml:"default_view" json:"default_view"`
	DefaultTheme string     `yaml:"default_theme" json:"default_theme"`
	PasswordHash string     `yaml:"password_hash,omitempty" json:"password_hash,omitempty"`
	CustomCss    string     `yaml:"custom_css,omitempty" json:"custom_css,omitempty"`
	PaletteHotkey string    `yaml:"palette_hotkey,omitempty" json:"palette_hotkey,omitempty"`
	Health       Health     `yaml:"health" json:"health"`
	Appearance   Appearance `yaml:"appearance" json:"appearance"`
	Sections     []Section  `yaml:"sections" json:"sections"`
	Links        []Link     `yaml:"links" json:"links"`
}

type Health struct {
	Enabled         bool `yaml:"enabled" json:"enabled"`
	IntervalSeconds int  `yaml:"interval_seconds" json:"interval_seconds"`
	TimeoutSeconds  int  `yaml:"timeout_seconds" json:"timeout_seconds"`
}

type Appearance struct {
	Page      Page      `yaml:"page" json:"page"`
	Grid      Grid      `yaml:"grid" json:"grid"`
	Item      Item      `yaml:"item" json:"item"`
	Icon      IconSizes `yaml:"icon" json:"icon"`
	Text      Text      `yaml:"text" json:"text"`
	StatusDot StatusDot `yaml:"status_dot" json:"status_dot"`
}

type Page struct {
	MaxWidth   int    `yaml:"max_width" json:"max_width"`
	Background string `yaml:"background" json:"background"`
	FontFamily string `yaml:"font_family" json:"font_family"`
	FontSize   int    `yaml:"font_size" json:"font_size"`
}

type Grid struct {
	Columns      string `yaml:"columns" json:"columns"`
	MinItemWidth int    `yaml:"min_item_width" json:"min_item_width"`
	Gap          int    `yaml:"gap" json:"gap"`
}

type Item struct {
	CornerRadius   int  `yaml:"corner_radius" json:"corner_radius"`
	Padding        int  `yaml:"padding" json:"padding"`
	Background     bool `yaml:"background" json:"background"`
	Border         bool `yaml:"border" json:"border"`
	BorderStrength int  `yaml:"border_strength" json:"border_strength"`
	Shadow         bool `yaml:"shadow" json:"shadow"`
	ShadowStrength int  `yaml:"shadow_strength" json:"shadow_strength"`
}

type IconSizes struct {
	SizeDefault int `yaml:"size_default" json:"size_default"`
	SizeCompact int `yaml:"size_compact" json:"size_compact"`
	SizeCard    int `yaml:"size_card" json:"size_card"`
	SizeLarge   int `yaml:"size_large" json:"size_large"`
}

type Text struct {
	Align           string `yaml:"align" json:"align"`
	ShowDescription bool   `yaml:"show_description" json:"show_description"`
}

type StatusDot struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Size     int    `yaml:"size" json:"size"`
	Position string `yaml:"position" json:"position"`
}

type Section struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

type Link struct {
	Type        string `yaml:"type,omitempty" json:"type,omitempty"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Text        string `yaml:"text,omitempty" json:"text,omitempty"`
	URL         string `yaml:"url" json:"url"`
	Icon        string `yaml:"icon" json:"icon"`
	Color       string `yaml:"color" json:"color"`
	Section     string `yaml:"section,omitempty" json:"section,omitempty"`
	Health      *bool  `yaml:"health,omitempty" json:"health,omitempty"`
}

// IsNote reports whether this item is a note (rather than a link).
func (l Link) IsNote() bool { return l.Type == "note" }

// Default returns a Config populated with sensible defaults.
func Default() Config {
	t := true
	return Config{
		Title:        "Minidash",
		DefaultView:  "default",
		DefaultTheme: "auto",
		PaletteHotkey: "ctrl+p",
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
	c.ApplyDefaults()
	return &c, nil
}

func (c *Config) ApplyDefaults() {
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
	if c.PaletteHotkey == "" {
		c.PaletteHotkey = d.PaletteHotkey
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
