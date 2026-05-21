package sonarqube

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestResourceProject_InternalValidate(t *testing.T) {
	t.Parallel()

	r := resourceProject()
	if err := r.InternalValidate(r.Schema, true); err != nil {
		t.Fatalf("resourceProject failed internal validation: %s", err)
	}
}

func TestResourceProject_Schema(t *testing.T) {
	t.Parallel()

	r := resourceProject()

	cases := []struct {
		field    string
		typ      schema.ValueType
		required bool
		optional bool
		computed bool
		def      interface{}
	}{
		{"key", schema.TypeString, true, false, false, nil},
		{"name", schema.TypeString, true, false, false, nil},
		{"public", schema.TypeBool, false, true, false, false},
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
			if s.Optional != tc.optional {
				t.Errorf("Optional = %v, want %v", s.Optional, tc.optional)
			}
			if s.Computed != tc.computed {
				t.Errorf("Computed = %v, want %v", s.Computed, tc.computed)
			}
			if tc.def != nil && s.Default != tc.def {
				t.Errorf("Default = %v, want %v", s.Default, tc.def)
			}
		})
	}
}

func TestResourceProject_CRUD(t *testing.T) {
	t.Parallel()

	r := resourceProject()

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
	if r.Exists == nil {
		t.Error("Exists is nil")
	}
	if r.Importer == nil {
		t.Error("Importer is nil; resource should support `terraform import`")
	}
}
