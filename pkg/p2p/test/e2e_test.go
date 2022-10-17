/*
   Copyright The Accelerated Container Image Authors

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

package test

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/data-accelerator/dadi-p2proxy/pkg/p2p/server"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setup() {
	rand.Seed(time.Now().UnixNano())
	log.SetLevel(log.InfoLevel)
	// ignore test.root
	flag.Bool("test.root", false, "")
}

func teardown() {
	if err := os.RemoveAll(server.Media); err != nil {
		fmt.Printf("Remove %s failed! %s", server.Media, err)
	}
}

func TestMain(m *testing.M) {
	log.Info("Start test!")
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestDockerPull(t *testing.T) {
	Assert := assert.New(t)
	serverList := server.StartServers(1, 2, true, true)
	for _, p2pServer := range serverList {
		port, _ := strconv.Atoi(p2pServer.Addr[1:])
		port += 100
		server.ConfigureDocker(fmt.Sprintf(":%d", port))
		_, code := server.ExecuteCmd(true, "docker", "pull", "wordpress")
		Assert.Equal(0, code)
		server.ExecuteCmd(true, "docker", "rmi", "-f", "wordpress")
		server.ExecuteCmd(true, "docker", "image", "prune", "-f")
	}
	for _, p2pServer := range serverList {
		_ = p2pServer.Shutdown(context.TODO())
	}
}
