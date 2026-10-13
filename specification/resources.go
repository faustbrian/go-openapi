// Package specification embeds the exact specification resources pinned by
// this module. Read returns copied bytes from embedded schema and IANA registry
// resources, not the complete normative Markdown corpus or a network loader.
// Names are embed.FS paths; unsupported names return an error. Authentication is
// performed by maintainer provenance gates against manifest pins, not afresh on
// each Read, and informative schemas do not supersede normative prose.
package specification

import (
	"embed"
	"fmt"
)

// resources contains the immutable files recorded by manifest.json.
//
//go:embed schemas registries/iana
var resources embed.FS

// Read returns a caller-owned copy of one pinned resource.
func Read(name string) ([]byte, error) {
	value, err := resources.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read pinned specification resource %q: %w", name, err)
	}
	return append([]byte(nil), value...), nil
}
