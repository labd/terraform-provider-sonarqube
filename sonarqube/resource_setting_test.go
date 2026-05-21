package sonarqube

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestResourceSettingsValue_InternalValidate(t *testing.T) {
	t.Parallel()

	r := resourceSettingsValue()
	if err := r.InternalValidate(r.Schema, true); err != nil {
		t.Fatalf("resourceSettingsValue failed internal validation: %s", err)
	}
}

func TestResourceSettingsValue_Schema(t *testing.T) {
	t.Parallel()

	r := resourceSettingsValue()

	cases := []struct {
		field    string
		typ      schema.ValueType
		required bool
	}{
		{"key", schema.TypeString, true},
		{"value", schema.TypeString, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			t.Parallel()

			s, ok := r.Schema[tc.field]
			if !ok {
				t.Fatalf("schema missing field %q", tc.field)
			}
			if s.Type != tc.typ {
				t.Errorf("Type = %v, want %v", s.Type, tc.typ)
			}
			if s.Required != tc.required {
				t.Errorf("Required = %v, want %v", s.Required, tc.required)
			}
		})
	}
}

func TestResourceSettingsValue_CRUD(t *testing.T) {
	t.Parallel()

	r := resourceSettingsValue()

	if r.Create == nil {
		t.Error("Create is nil")
	}
	if r.Read == nil {
		t.Error("Read is nil")
	}
	if r.Update == nil {
		t.Error("Update is nil")
	}
	if r.Delete == nil {
		t.Error("Delete is nil")
	}
	if r.Importer == nil {
		t.Error("Importer is nil; resource should support `terraform import`")
	}
}
