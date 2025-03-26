/*
Copyright 2025 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg_test

import (
	"testing"

	"github.com/linuxsuren/atest-ext-store-iotdb/pkg"
	"github.com/stretchr/testify/assert"
)

func TestInnerSQL(t *testing.T) {
	t.Run("support postgres", func(t *testing.T) {
		assert.NotNil(t, pkg.GetInnerSQL(pkg.DialectorPostgres))
	})

	t.Run("complext native SQL", func(t *testing.T) {
		const sql = "select * from table limit 3"
		assert.Equal(t, sql, pkg.GetInnerSQL(pkg.DialectorMySQL).ToNativeSQL(sql))
		assert.Equal(t, sql, pkg.GetInnerSQL(pkg.DialectorPostgres).ToNativeSQL(sql))
	})

	innerSQLs := []string{
		pkg.InnerSelectTable_,
		pkg.InnerSelectTableLimit_,
		pkg.InnerShowDatabases,
		pkg.InnerShowTables,
		pkg.InnerCurrentDB,
	}
	t.Run("inner SQL", func(t *testing.T) {
		for _, sql := range innerSQLs {
			assert.NotEqual(t, sql, pkg.GetInnerSQL(pkg.DialectorMySQL).ToNativeSQL(sql))
			assert.NotEqual(t, sql, pkg.GetInnerSQL(pkg.DialectorPostgres).ToNativeSQL(sql))
		}
	})
}
