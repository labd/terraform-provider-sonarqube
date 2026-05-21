package sonarqube

import (
	"context"
	"testing"
	"time"

	sonargo "github.com/labd/sonargo/sonar"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tfplugin "github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	t.Parallel()

	var _ terraform.ResourceProvider = Provider()
}

func TestProvider_Schema(t *testing.T) {
	// No t.Parallel(): subtests use t.Setenv to verify the EnvDefaultFunc,
	// which is incompatible with parallel execution at any level.

	p := Provider().(*schema.Provider)

	cases := []struct {
		field     string
		typ       schema.ValueType
		required  bool
		sensitive bool
		envVar    string
	}{
		{"url", schema.TypeString, true, true, "SONAR_URL"},
		{"token", schema.TypeString, true, true, "SONAR_TOKEN"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			// No t.Parallel(): t.Setenv is incompatible with parallel subtests.

			s, ok := p.Schema[tc.field]
			if !ok {
				t.Fatalf("provider schema missing field %q", tc.field)
			}
			if s.Type != tc.typ {
				t.Errorf("Type = %v, want %v", s.Type, tc.typ)
			}
			if s.Required != tc.required {
				t.Errorf("Required = %v, want %v", s.Required, tc.required)
			}
			if s.Sensitive != tc.sensitive {
				t.Errorf("Sensitive = %v, want %v", s.Sensitive, tc.sensitive)
			}
			if s.DefaultFunc == nil {
				t.Errorf("DefaultFunc is nil; expected EnvDefaultFunc(%q, ...)", tc.envVar)
			} else {
				t.Setenv(tc.envVar, "sentinel-value")
				v, err := s.DefaultFunc()
				if err != nil {
					t.Errorf("DefaultFunc returned error: %s", err)
				}
				if v != "sentinel-value" {
					t.Errorf("DefaultFunc did not read from %s; got %v", tc.envVar, v)
				}
			}
		})
	}
}

func TestProvider_RegistersResources(t *testing.T) {
	t.Parallel()

	p := Provider().(*schema.Provider)

	expected := []string{
		"sonarqube_project",
		"sonarqube_project_user",
		"sonarqube_settings_value",
	}
	for _, name := range expected {
		if _, ok := p.ResourcesMap[name]; !ok {
			t.Errorf("resource %q not registered on provider", name)
		}
	}
}

func TestProvider_RegistersDataSources(t *testing.T) {
	t.Parallel()

	p := Provider().(*schema.Provider)

	expected := []string{"sonarqube_user"}
	for _, name := range expected {
		if _, ok := p.DataSourcesMap[name]; !ok {
			t.Errorf("data source %q not registered on provider", name)
		}
	}
}

// TestProvider_ConfigureFunc_RoundTrip verifies that providerConfigure returns
// a non-nil *sonargo.Client when given valid url+token schema values.
func TestProvider_ConfigureFunc_RoundTrip(t *testing.T) {
	t.Parallel()

	p := Provider().(*schema.Provider)
	if p.ConfigureFunc == nil {
		t.Fatal("ConfigureFunc is nil")
	}

	d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{
		"url":   "http://sonarqube.example.test",
		"token": "test-token",
	})

	meta, err := p.ConfigureFunc(d)
	if err != nil {
		t.Fatalf("ConfigureFunc returned error: %s", err)
	}
	if meta == nil {
		t.Fatal("ConfigureFunc returned nil meta; expected *sonargo.Client")
	}
	if _, ok := meta.(*sonargo.Client); !ok {
		t.Fatalf("ConfigureFunc returned %T; want *sonargo.Client", meta)
	}
}

// TestProvider_PluginServe verifies the go-plugin / gRPC handshake by starting
// the provider via DebugServe and confirming a usable ReattachConfig is
// produced. It does not drive an RPC call; a successful handshake is enough
// to catch breakage in the plugin transport stack.
func TestProvider_PluginServe(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reattach, closeCh, err := tfplugin.DebugServe(ctx, &tfplugin.ServeOpts{
		ProviderFunc: Provider,
	})
	if err != nil {
		t.Fatalf("DebugServe returned error: %s", err)
	}

	if reattach.Pid == 0 {
		t.Error("ReattachConfig.Pid is 0; plugin server did not start cleanly")
	}
	if reattach.Protocol == "" {
		t.Error("ReattachConfig.Protocol is empty; handshake did not complete")
	}
	if reattach.Addr.Network == "" || reattach.Addr.String == "" {
		t.Errorf("ReattachConfig.Addr is incomplete: %+v", reattach.Addr)
	}

	// Tear the server down and confirm it shuts down within the deadline.
	cancel()
	select {
	case <-closeCh:
	case <-time.After(5 * time.Second):
		t.Fatal("plugin server did not shut down after context cancel")
	}
}
