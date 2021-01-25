package azure

import (
	"testing"

	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
)

func TestAccAzureRMValidateKeyVaultChildID(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := keyVaultValidate.NestedItemId(tc.Input, "example")
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			return
		}

		if tc.ExpectError && len(warnings) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(warnings) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}

func TestAccAzureRMValidateKeyVaultChildIDVersionOptional(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := keyVaultValidate.NestedItemIdWithOptionalVersion(tc.Input, "example")
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			return
		}

		if tc.ExpectError && len(warnings) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(warnings) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}

func TestAccAzureRMKeyVaultChild_parseID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    keyVaultParse.NestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "hello",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "world",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "1492",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		secretId, err := keyVaultParse.ParseNestedItemID(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
			}

			return
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}
	}
}

func TestAccAzureRMKeyVaultChild_parseIDVersionOptional(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    keyVaultParse.NestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "hello",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "world",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: keyVaultParse.NestedItemId{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "1492",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		secretId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
			}

			return
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}
	}
}
