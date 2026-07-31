package model

import (
	"strconv"

	"minidash/internal/config"
)

// Links are referenced by stringified positional index into c.Links.
// Sufficient for reorder; keeps YAML human-friendly (no UUIDs).

func LinkIDsByOrder(c *config.Config) []string {
	out := make([]string, len(c.Links))
	for i := range c.Links {
		out[i] = strconv.Itoa(i)
	}
	return out
}

func ReorderLinksByIDs(c *config.Config, ids []string) {
	used := map[int]bool{}
	reordered := make([]config.Link, 0, len(c.Links))
	for _, id := range ids {
		n, err := strconv.Atoi(id)
		if err != nil || n < 0 || n >= len(c.Links) || used[n] {
			continue
		}
		used[n] = true
		reordered = append(reordered, c.Links[n])
	}
	for i := range c.Links {
		if !used[i] {
			reordered = append(reordered, c.Links[i])
		}
	}
	c.Links = reordered
}

func ReorderSectionsByIDs(c *config.Config, ids []string) {
	pos := map[string]int{}
	for i, id := range ids {
		if _, ok := pos[id]; !ok {
			pos[id] = i
		}
	}
	for i := 0; i < len(c.Sections); i++ {
		for j := i + 1; j < len(c.Sections); j++ {
			pi, pj := len(c.Sections), len(c.Sections)
			if v, ok := pos[c.Sections[i].ID]; ok {
				pi = v
			}
			if v, ok := pos[c.Sections[j].ID]; ok {
				pj = v
			}
			if pj < pi {
				c.Sections[i], c.Sections[j] = c.Sections[j], c.Sections[i]
			}
		}
	}
}

func NextSectionID(c config.Config) string {
	taken := map[string]bool{}
	for _, s := range c.Sections {
		taken[s.ID] = true
	}
	for i := 1; ; i++ {
		id := "s" + strconv.Itoa(i)
		if !taken[id] {
			return id
		}
	}
}
