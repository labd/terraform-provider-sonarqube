package sonarqube

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestResourceProjectUser_InternalValidate(t *testing.T) {
	t.Parallel()

	r := resourceProjectUser()
	if err := r.InternalValidate(r.Schema, true); err != nil {
		t.Fatalf("resourceProjectUser failed internal validation: %s", err)
	}
}

func TestResourceProjectUser_Schema(t *testing.T) {
	t.Parallel()

	r := resourceProjectUser()

	cases := []struct {
		field    string
		typ      schema.ValueType
		required bool
		forceNew bool
	}{
		{"login", schema.TypeString, true, true},
		{"project_key", schema.TypeString, true, true},
		{"permissions", schema.TypeList, true, true},
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
			if s.ForceNew != tc.forceNew {
				t.Errorf("ForceNew = %v, want %v", s.ForceNew, tc.forceNew)
			}
		})
	}
}

func TestResourceProjectUser_CRUD(t *testing.T) {
	t.Parallel()

	r := resourceProjectUser()

	if r.Create == nil {
		t.Error("Create is nil")
	}
	if r.Read == nil {
		t.Error("Read is nil")
	}
	if r.Delete == nil {
		t.Error("Delete is nil")
	}
	// Update is intentionally omitted: every attribute is ForceNew, so changes
	// are handled via destroy + recreate.
	if r.Update != nil {
		t.Error("Update should be nil; all fields are ForceNew")
	}
}

func TestResourceProjectUser_PermissionsValidation(t *testing.T) {
	t.Parallel()

	r := resourceProjectUser()
	elem, ok := r.Schema["permissions"].Elem.(*schema.Schema)
	if !ok {
		t.Fatalf("permissions Elem is not *schema.Schema (got %T)", r.Schema["permissions"].Elem)
	}
	if elem.ValidateFunc == nil {
		t.Fatal("permissions Elem.ValidateFunc is nil")
	}

	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid_admin", "admin", false},
		{"valid_codeviewer", "codeviewer", false},
		{"valid_issueadmin", "issueadmin", false},
		{"valid_securityhotspotadmin", "securityhotspotadmin", false},
		{"valid_scan", "scan", false},
		{"valid_user", "user", false},
		{"invalid_empty", "", true},
		{"invalid_uppercase", "ADMIN", true},
		{"invalid_unknown", "owner", true},
		{"invalid_typo", "admins", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, errs := elem.ValidateFunc(tc.input, "permissions")
			gotErr := len(errs) > 0
			if gotErr != tc.wantErr {
				t.Errorf("ValidateFunc(%q): gotErr=%v wantErr=%v errs=%v",
					tc.input, gotErr, tc.wantErr, errs)
			}
		})
	}
}
