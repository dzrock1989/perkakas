package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	type omama struct {
		Test  string `filter:"test"`
		Test2 string `filter:"test2"`
		Name  string `filter:"name"`
	}

	var q string
	var a []interface{}
	var e error

	q, a, e = BuildFilterQuery(omama{}, "test:like:value")
	require.NoError(t, e)
	require.Equal(t, "omamas.test ILIKE ?", q)
	require.Equal(t, "%value%", a[0].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:neq:value")
	require.NoError(t, e)
	require.Equal(t, "omamas.test != ?", q)
	require.Equal(t, "value", a[0].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:value")
	require.NoError(t, e)
	require.Equal(t, "omamas.test = ?", q)
	require.Equal(t, "value", a[0].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:1")
	require.NoError(t, e)
	require.Equal(t, "omamas.test = ?", q)
	require.Equal(t, "1", a[0].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:in:1,2,3")
	require.NoError(t, e)
	require.Equal(t, "omamas.test IN ?", q)
	require.Equal(t, 3, len(a[0].([]string)))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:null")
	require.NoError(t, e)
	require.Equal(t, "omamas.test is null", q)
	require.Equal(t, 0, len(a))

	q, a, e = BuildFilterQuery(omama{}, "test:in:a,b,c")
	require.NoError(t, e)
	require.Equal(t, "omamas.test IN ?", q)
	require.Equal(t, 3, len(a[0].([]string)))

	// _, _, e = BuildFilterQuery(omama{}, "test:in:1,2,3,a")
	// require.EqualError(t, e, "invalid values")

	_, _, e = BuildFilterQuery(omama{}, "xxx:eq:value")
	require.EqualError(t, e, "invalid query filter for 'xxx'")

	_, _, e = BuildFilterQuery(omama{}, "test::value")
	require.EqualError(t, e, "invalid clause")

	_, _, e = BuildFilterQuery(omama{}, "test:value")
	require.EqualError(t, e, "invalid query filter")

	_, _, e = BuildFilterQuery(omama{}, "test:omama:value")
	require.EqualError(t, e, "invalid clause")

	_, _, e = BuildFilterQuery(omama{}, "-:omama:value")
	require.EqualError(t, e, "invalid query filter for '-'")

	q, a, e = BuildFilterQuery(omama{}, "test:startWith:value")
	require.NoError(t, e)
	require.Equal(t, "omamas.test ILIKE ?", q)
	require.Equal(t, "value%", a[0].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:endWith:value")
	require.NoError(t, e)
	require.Equal(t, "omamas.test ILIKE ?", q)
	require.Equal(t, "%value", a[0].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:value;OR;name:eq:myName")
	require.NoError(t, e)
	require.Equal(t, "omamas.test = ? OR omamas.name = ?", q)
	require.Equal(t, "value", a[0].(string))
	require.Equal(t, "myName", a[1].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:value;AND;name:eq:myName")
	require.NoError(t, e)
	require.Equal(t, "omamas.test = ? AND omamas.name = ?", q)
	require.Equal(t, "value", a[0].(string))
	require.Equal(t, "myName", a[1].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:value;AND;name:eq:ANDA")
	require.NoError(t, e)
	require.Equal(t, "omamas.test = ? AND omamas.name = ?", q)
	require.Equal(t, "value", a[0].(string))
	require.Equal(t, "ANDA", a[1].(string))

	q, a, e = BuildFilterQuery(omama{}, "test:eq:value;OR;name:eq:ORMAS")
	require.NoError(t, e)
	require.Equal(t, "omamas.test = ? OR omamas.name = ?", q)
	require.Equal(t, "value", a[0].(string))
	require.Equal(t, "ORMAS", a[1].(string))

	_, _, e = BuildFilterQuery(omama{}, "test:omama:value;;name:eq:myName")
	require.EqualError(t, e, "invalid query filter")
}
