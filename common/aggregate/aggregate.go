package aggregate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tigapilarmandiri/perkakas/common/middlewares/authorization"
)

var (
	errJenisWilayahNotValid = errors.New("jenis wilayah tidak valid")
	errNotPermitted         = errors.New("you're not permitted")
)

var mapKepolisian = map[string]int{
	"mabes":  0,
	"polda":  1,
	"polres": 2,
	"polsek": 3,
}

func GetViewAggregateName(tableName, jenisWilayah string, claims authorization.Claims) (string, error) {
	var viewName string
	kepolisian := strings.ToLower(claims.KepolisianLevel)

	switch jenisWilayah {
	case "kl":
		viewName = fmt.Sprintf("v_%s_desa_lurah", tableName)
	case "kc":
		viewName = fmt.Sprintf("v_%s_kecamatan", tableName)
	case "kt":
		viewName = fmt.Sprintf("v_%s_kabkota", tableName)
		if i, ok := mapKepolisian[kepolisian]; ok && !claims.IsSuperadmin {
			if i > 2 {
				return "", errNotPermitted
			}
		}
	case "pr":
		viewName = fmt.Sprintf("v_%s_provinsi", tableName)
		if i, ok := mapKepolisian[kepolisian]; ok && !claims.IsSuperadmin {
			if i > 1 {
				return "", errNotPermitted
			}
		}
	case "ns":
		viewName = fmt.Sprintf("v_%s_nasional", tableName)
		if i, ok := mapKepolisian[kepolisian]; ok && !claims.IsSuperadmin {
			if i > 0 {
				return "", errNotPermitted
			}
		}
	default:
		return "", errJenisWilayahNotValid
	}

	return viewName, nil
}

func GetViewAggregateNameDetail(tableName, jenisWilayah string, claims authorization.Claims) (string, error) {
	var viewName string
	kepolisian := strings.ToLower(claims.KepolisianLevel)

	switch jenisWilayah {
	case "kl":
		viewName = fmt.Sprintf("v_%s_desa_lurah_detail", tableName)
	case "kc":
		viewName = fmt.Sprintf("v_%s_kecamatan_detail", tableName)
	case "kt":
		viewName = fmt.Sprintf("v_%s_kabkota_detail", tableName)
		if i, ok := mapKepolisian[kepolisian]; ok && !claims.IsSuperadmin {
			if i > 2 {
				return "", errNotPermitted
			}
		}
	case "pr":
		viewName = fmt.Sprintf("v_%s_provinsi_detail", tableName)
		if i, ok := mapKepolisian[kepolisian]; ok && !claims.IsSuperadmin {
			if i > 1 {
				return "", errNotPermitted
			}
		}
	case "ns":
		viewName = fmt.Sprintf("v_%s_nasional_detail", tableName)
		if i, ok := mapKepolisian[kepolisian]; ok && !claims.IsSuperadmin {
			if i > 0 {
				return "", errNotPermitted
			}
		}
	default:
		return "", errJenisWilayahNotValid
	}

	return viewName, nil
}
