package config

import "github.com/derailed/k9s/internal/client"

// CustomResourceLink tracks K9s CustomResourceLink configuration.
type CustomResourceLink struct {
	Target        string            `yaml:"target"`
	LabelSelector map[string]string `yaml:"labelSelector,omitempty"`
	FieldSelector map[string]string `yaml:"fieldSelector,omitempty"`
}

// NewResourceLink creates a new CustomResourceLink configuration.
func NewResourceLink() *CustomResourceLink {
	return &CustomResourceLink{}
}

// Validate a CustomResourceLink config.
func (c *CustomResourceLink) Validate(conn client.Connection, ks KubeSettings) {

}
