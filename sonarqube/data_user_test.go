package sonarqube

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestDataSourceUser_InternalValidate(t *testing.T) {
	t.Parallel()

	r := dataSourceUser()
	if err := r.InternalValidate(r.Schema, false); err != nil {
		t.Fatalf("dataSourceUser failed internal validation: %s", err)
	}
}

func TestDataSourceUser_Schema(t *testing.T) {
	t.Parallel()

	r := dataSourceUser()

	cases := []struct {
		field    string
		typ      schema.ValueType
		required bool
		computed bool
	}{
		{"email", schema.TypeString, true, false},
		{"login", schema.TypeString, false, true},
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
			if s.Computed != tc.computed {
				t.Errorf("Computed = %v, want %v", s.Computed, tc.computed)
			}
		})
	}
}

func TestDataSourceUser_OnlyRead(t *testing.T) {
	t.Parallel()

	r := dataSourceUser()

	if r.Read == nil {
		t.Error("Read is nil")
	}
	// Data sources must not declare write CRUD operations.
	if r.Create != nil {
		t.Error("data source should not have Create")
	}
	if r.Update != nil {
		t.Error("data source should not have Update")
	}
	if r.Delete != nil {
		t.Error("data source should not have Delete")
	}
}
