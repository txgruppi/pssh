package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

type Host struct {
	Name         string   `json:"name"`
	Addr         string   `json:"addr"`
	User         string   `json:"user"`
	SudoPassword string   `json:"sudo_password"`
	Tags         []string `json:"tags"`
}

type Hosts []Host

func (t Hosts) FilterByTag(tag string) Hosts {
	res := make(Hosts, 0, len(t))
	for _, h := range t {
		if slices.Contains(h.Tags, tag) {
			res = append(res, h)
		}
	}
	return res
}

func (t Hosts) FilterByTags(tags []string) Hosts {
	res := t
	for _, tag := range tags {
		res = res.FilterByTag(tag)
	}
	return res
}

func loadHostsFromFile(filepath string) (Hosts, error) {
	var hosts Hosts
	file, err := os.Open(filepath)
	if err != nil {
		return hosts, fmt.Errorf("error opening hosts file %s: %w", filepath, err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&hosts); err != nil {
		return hosts, fmt.Errorf("error decoding hosts file %s: %w", filepath, err)
	}
	return hosts, nil
}
