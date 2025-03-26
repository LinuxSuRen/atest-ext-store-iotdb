/*
Copyright 2023-2025 API Testing Authors.

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
	"time"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/atest-ext-store-iotdb/pkg"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConvertToRemoteTestCase(t *testing.T) {
	result := pkg.ConvertToRemoteTestCase(&pkg.TestCase{
		Name:   "name",
		API:    "api",
		Method: "get",
		Body:   "body",
		Query:  sampleJSONMap,
		Header: sampleJSONMap,
		Form:   sampleJSONMap,

		ExpectStatusCode: 1,
		ExpectBody:       "expect body",
		ExpectSchema:     "schema",
		ExpectHeader:     sampleJSONMap,
		ExpectBodyFields: sampleJSONMap,
		ExpectVerify:     `["one"]`,
	})
	assert.Equal(t, &server.TestCase{
		Name: "name",
		Request: &server.Request{
			Api:    "api",
			Method: "get",
			Body:   "body",
			Query:  samplePairs,
			Header: samplePairs,
			Form:   samplePairs,
		},
		Response: &server.Response{
			StatusCode:       1,
			BodyFieldsExpect: samplePairs,
			Verify:           []string{"one"},
			Schema:           "schema",
			Body:             "expect body",
			Header:           samplePairs,
		},
	}, result)
}

func TestConverToDBTestCase(t *testing.T) {
	t.Run("without request and response", func(t *testing.T) {
		result := pkg.ConverToDBTestCase(&server.TestCase{})
		assert.Equal(t, &pkg.TestCase{}, result)
	})

	t.Run("only have request", func(t *testing.T) {
		result := pkg.ConverToDBTestCase(&server.TestCase{
			Request: &server.Request{
				Api:    "api",
				Method: "get",
				Body:   "body",
				Header: samplePairs,
				Cookie: samplePairs,
				Query:  samplePairs,
				Form:   samplePairs,
			},
		})
		assert.Equal(t, &pkg.TestCase{
			API:    "api",
			Method: "get",
			Body:   "body",
			Header: sampleJSONMap,
			Cookie: sampleJSONMap,
			Query:  sampleJSONMap,
			Form:   sampleJSONMap,
		}, result)
	})

	t.Run("only have response", func(t *testing.T) {
		result := pkg.ConverToDBTestCase(&server.TestCase{
			Response: &server.Response{
				StatusCode:       1,
				Body:             "body",
				Schema:           "schema",
				Header:           samplePairs,
				BodyFieldsExpect: samplePairs,
				Verify:           []string{"one"},
			},
		})
		assert.Equal(t, &pkg.TestCase{
			ExpectBody:       "body",
			ExpectStatusCode: 1,
			ExpectSchema:     "schema",
			ExpectVerify:     `["one"]`,
			ExpectHeader:     sampleJSONMap,
			ExpectBodyFields: sampleJSONMap,
		}, result)
	})
}

func TestConvertTestSuite(t *testing.T) {
	t.Run("ConvertToDBTestSuite", func(t *testing.T) {
		result := pkg.ConvertToDBTestSuite(&remote.TestSuite{
			Name:  "name",
			Api:   "api",
			Param: samplePairs,
			Spec: &server.APISpec{
				Kind: "kind",
			},
		})
		assert.Equal(t, &pkg.TestSuite{
			Name:     "name",
			API:      "api",
			SpecKind: "kind",
			Param:    `{"key":"value"}`,
		}, result)
	})

	t.Run("ConvertToGRPCTestSuite", func(t *testing.T) {
		result := pkg.ConvertToGRPCTestSuite(&pkg.TestSuite{
			Name: "name",
			API:  "api",
		})
		assert.Equal(t, &remote.TestSuite{
			Name: "name",
			Api:  "api",
			Spec: &server.APISpec{},
		}, result)
	})

	t.Run("sliceToJSON", func(t *testing.T) {
		assert.Equal(t, "[]", pkg.SliceToJSON(nil))
	})
}

func TestConvertToDBHistoryTestResult(t *testing.T) {
	t.Run("without testcaseResult and historyTestcase", func(t *testing.T) {
		result := pkg.ConvertToDBHistoryTestResult(&server.HistoryTestResult{})
		assert.Equal(t, &pkg.HistoryTestResult{}, result)
	})

	t.Run("have testcaseResult", func(t *testing.T) {
		result := pkg.ConvertToDBHistoryTestResult(&server.HistoryTestResult{
			Data: &server.HistoryTestCase{
				SuiteParam: []*server.Pair{{
					Key:   "key",
					Value: "value",
				}},
				Request: &server.Request{
					Header: []*server.Pair{{
						Key:   "key",
						Value: "value",
					}},
				},
				Response: &server.Response{
					Header: []*server.Pair{{
						Key:   "key",
						Value: "value",
					}},
				},
				SuiteSpec: &server.APISpec{
					Kind: "kind",
				},
			},
			TestCaseResult: []*server.TestCaseResult{
				{
					StatusCode: 200,
					Body:       "Test body",
					Output:     "Test output",
				},
			},
		})
		assert.Equal(t, &pkg.HistoryTestResult{
			StatusCode:       200,
			Body:             "Test body",
			Output:           "Test output",
			Param:            `{"key":"value"}`,
			Header:           `{"key":"value"}`,
			ExpectHeader:     `{"key":"value"}`,
			ExpectBodyFields: "{}",
			ExpectVerify:     "[]",
			Message:          "",
			ExpectSchema:     "",
			HistoryHeader:    `{}`,
			Cookie:           `{}`,
			Query:            `{}`,
			Form:             `{}`,
			SpecKind:         "kind",
			SpecURL:          "",
		}, result)
	})
}

var now = time.Now().UTC()
var nowString = now.Format("2006-01-02T15:04:05.999999999")

func TestConvertToRemoteHistoryTestResult(t *testing.T) {
	assert.Equal(t, &server.HistoryTestResult{
		CreateTime: timestamppb.New(now),
		TestCaseResult: []*server.TestCaseResult{
			{
				Body:   "body",
				Output: "output",
				Header: samplePairs,
			},
		},
		Data: &server.HistoryTestCase{
			CreateTime: timestamppb.New(now),
			SuiteSpec: &server.APISpec{
				Kind: "kind",
			},
			Request: &server.Request{
				Api:    "",
				Method: "",
				Header: []*server.Pair{
					{Key: "key", Value: "value"},
				},
				Body: "body",
			},
			Response: &server.Response{
				StatusCode: 0,
				Body:       "",
				Header:     nil,
			},
		},
	}, pkg.ConvertToRemoteHistoryTestResult(&pkg.HistoryTestResult{
		CreateTime: nowString,
		Body:       "body",
		Output:     "output",
		Header:     sampleJSONMap,
		SpecKind:   "kind",
	}))
}

func TestConvertToGRPCHistoryTestSuite(t *testing.T) {
	assert.Equal(t, &remote.HistoryTestSuite{
		Items: []*server.HistoryTestCase{
			{
				CreateTime: timestamppb.New(now),
				SuiteName:  "name",
				SuiteSpec: &server.APISpec{
					Kind: "kind",
				},
				Request: &server.Request{
					Body: "Test Body",
				},
				Response: &server.Response{},
			},
		},
	}, pkg.ConvertToGRPCHistoryTestSuite(&pkg.HistoryTestResult{
		CreateTime: nowString,
		SuiteName:  "name",
		Body:       "Test Body",
		SpecKind:   "kind",
	}))
}

func TestConvertToGRPCHistoryTestCase(t *testing.T) {
	assert.Equal(t, &server.HistoryTestCase{
		CreateTime: timestamppb.New(now),
		SuiteName:  "name",
		SuiteSpec: &server.APISpec{
			Kind: "kind",
		},
		Request: &server.Request{
			Body: "Test Body",
		},
		Response: &server.Response{},
	}, pkg.ConvertToGRPCHistoryTestCase(&pkg.HistoryTestResult{
		CreateTime: nowString,
		SuiteName:  "name",
		Body:       "Test Body",
		SpecKind:   "kind",
	}))
}

const sampleJSONMap = `{"key":"value"}`

var samplePairs []*server.Pair = []*server.Pair{{
	Key:   "key",
	Value: "value",
}}
