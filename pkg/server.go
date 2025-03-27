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
package pkg

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/apache/iotdb-client-go/client"
	"github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/version"
)

type dbserver struct {
	remote.UnimplementedLoaderServer
	defaultHistoryLimit int
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer(defaultHistoryLimit int) (s remote.LoaderServer) {
	s = &dbserver{
		defaultHistoryLimit: defaultHistoryLimit,
	}
	return
}

func (s *dbserver) getClientWithDatabase(ctx context.Context) (dbQuery DataQuery, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to database")
	} else {
		port := "6667"
		host := store.URL
		if strings.Contains(store.URL, ":") {
			obj := strings.Split(store.URL, ":")
			if len(obj) > 1 {
				port = obj[1]
				host = obj[0]
			}
		}

		timezone := store.Properties["timezone"]
		config := &client.PoolConfig{
			Host:     host,
			Port:     port,
			UserName: store.Username,
			Password: store.Password,
			TimeZone: timezone,
		}
		sessionPool := client.NewSessionPool(config, 3, 60000, 60000, false)

		dbQuery = NewCommonDataQuery(GetInnerSQL("iotdb"), &sessionPool)
	}
	return
}

func (s *dbserver) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	return
}

func (s *dbserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	return
}

func (s *dbserver) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	return
}

func (s *dbserver) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *server.TestCases, err error) {
	return
}

func (s *dbserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) CreateTestCaseHistory(ctx context.Context, historyTestResult *server.HistoryTestResult) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) ListHistoryTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.HistoryTestSuites, err error) {
	return
}

func (s *dbserver) GetTestCase(ctx context.Context, testcase *server.TestCase) (result *server.TestCase, err error) {
	return
}

func (s *dbserver) GetHistoryTestCaseWithResult(ctx context.Context, testcase *server.HistoryTestCase) (result *server.HistoryTestResult, err error) {
	return
}

func (s *dbserver) GetHistoryTestCase(ctx context.Context, testcase *server.HistoryTestCase) (result *server.HistoryTestCase, err error) {
	return
}

func (s *dbserver) GetTestCaseAllHistory(ctx context.Context, testcase *server.TestCase) (result *server.HistoryTestCases, err error) {
	return
}

func (s *dbserver) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	return
}

func (s *dbserver) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) DeleteHistoryTestCase(ctx context.Context, historyTestCase *server.HistoryTestCase) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) DeleteAllHistoryTestCase(ctx context.Context, historyTestCase *server.HistoryTestCase) (reply *server.Empty, err error) {
	return
}

func (s *dbserver) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	reply = &server.ExtensionStatus{
		Ready: true,
	}
	return
}

func (s *dbserver) GetVersion(context.Context, *server.Empty) (ver *server.Version, err error) {
	ver = &server.Version{
		Version: version.GetVersion(),
		Commit:  version.GetCommit(),
		Date:    version.GetDate(),
	}
	return
}

func (s *dbserver) PProf(ctx context.Context, in *server.PProfRequest) (data *server.PProfData, err error) {
	log.Println("pprof", in.Name)

	data = &server.PProfData{
		Data: extension.LoadPProf(in.Name),
	}
	return
}
