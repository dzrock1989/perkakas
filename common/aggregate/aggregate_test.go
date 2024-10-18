package aggregate

import (
	"testing"

	"github.com/tigapilarmandiri/perkakas/common/middlewares/authorization"
)

func TestGetViewAggregateName(t *testing.T) {
	tests := []struct {
		name          string
		expected      string
		expectedError error
		tableName     string
		jenisWilayah  string
		claims        authorization.Claims
	}{
		{"success_polsek_desa", "v_data_desa_lurah", nil, "data", "kl", authorization.Claims{KepolisianLevel: "POLSEK"}},
		{"success_polsek_kecamatan", "v_data_kecamatan", nil, "data", "kc", authorization.Claims{KepolisianLevel: "POLSEK"}},
		{"success_polsek_kota", "", errNotPermitted, "data", "kt", authorization.Claims{KepolisianLevel: "POLSEK"}},
		{"success_polsek_provinsi", "", errNotPermitted, "data", "pr", authorization.Claims{KepolisianLevel: "POLSEK"}},
		{"success_polsek_nasional", "", errNotPermitted, "data", "ns", authorization.Claims{KepolisianLevel: "POLSEK"}},

		{"success_polres_desa", "v_data_desa_lurah", nil, "data", "kl", authorization.Claims{KepolisianLevel: "POLRES"}},
		{"success_polres_kecamatan", "v_data_kecamatan", nil, "data", "kc", authorization.Claims{KepolisianLevel: "POLRES"}},
		{"success_polres_kota", "v_data_kabkota", nil, "data", "kt", authorization.Claims{KepolisianLevel: "POLRES"}},
		{"success_polres_provinsi", "", errNotPermitted, "data", "pr", authorization.Claims{KepolisianLevel: "POLRES"}},
		{"success_polres_nasional", "", errNotPermitted, "data", "ns", authorization.Claims{KepolisianLevel: "POLRES"}},

		{"success_polda_desa", "v_data_desa_lurah", nil, "data", "kl", authorization.Claims{KepolisianLevel: "POLDA"}},
		{"success_polda_kecamatan", "v_data_kecamatan", nil, "data", "kc", authorization.Claims{KepolisianLevel: "POLDA"}},
		{"success_polda_kota", "v_data_kabkota", nil, "data", "kt", authorization.Claims{KepolisianLevel: "POLDA"}},
		{"success_polda_provinsi", "v_data_provinsi", nil, "data", "pr", authorization.Claims{KepolisianLevel: "POLDA"}},
		{"success_polda_nasional", "", errNotPermitted, "data", "ns", authorization.Claims{KepolisianLevel: "POLDA"}},

		{"success_mabes_desa", "v_data_desa_lurah", nil, "data", "kl", authorization.Claims{KepolisianLevel: "MABES"}},
		{"success_mabes_kecamatan", "v_data_kecamatan", nil, "data", "kc", authorization.Claims{KepolisianLevel: "MABES"}},
		{"success_mabes_kota", "v_data_kabkota", nil, "data", "kt", authorization.Claims{KepolisianLevel: "MABES"}},
		{"success_mabes_provinsi", "v_data_provinsi", nil, "data", "pr", authorization.Claims{KepolisianLevel: "MABES"}},
		{"success_mabes_nasional", "v_data_nasional", nil, "data", "ns", authorization.Claims{KepolisianLevel: "MABES"}},

		{"success_superadmin_desa", "v_data_desa_lurah", nil, "data", "kl", authorization.Claims{KepolisianLevel: "POLSEK", IsSuperadmin: true}},
		{"success_superadmin_kecamatan", "v_data_kecamatan", nil, "data", "kc", authorization.Claims{KepolisianLevel: "POLSEK", IsSuperadmin: true}},
		{"success_superadmin_kota", "v_data_kabkota", nil, "data", "kt", authorization.Claims{KepolisianLevel: "POLSEK", IsSuperadmin: true}},
		{"success_superadmin_provinsi", "v_data_provinsi", nil, "data", "pr", authorization.Claims{KepolisianLevel: "POLSEK", IsSuperadmin: true}},
		{"success_superadmin_nasional", "v_data_nasional", nil, "data", "ns", authorization.Claims{KepolisianLevel: "POLSEK", IsSuperadmin: true}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := GetViewAggregateName(tt.tableName, tt.jenisWilayah, tt.claims)
			if actual != tt.expected || err != tt.expectedError {
				t.Errorf("(%s, %s, %+v): expected (%s, %+v), actual (%s, %+v)",
					tt.tableName,
					tt.jenisWilayah,
					tt.claims,
					tt.expected,
					tt.expectedError,
					actual,
					err,
				)
			}
		})
	}
}
