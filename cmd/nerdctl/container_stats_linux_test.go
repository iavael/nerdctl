/*
   Copyright The containerd Authors.

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

package main

import (
	"fmt"
	"testing"

	"github.com/containerd/nerdctl/pkg/infoutil"
	"github.com/containerd/nerdctl/pkg/rootlessutil"
	"github.com/containerd/nerdctl/pkg/testutil"
)

func TestStats(t *testing.T) {
	// this comment is for `nerdctl ps` but it also valid for `nerdctl stats` :
	// https://github.com/containerd/nerdctl/pull/223#issuecomment-851395178
	if rootlessutil.IsRootless() && infoutil.CgroupsVersion() == "1" {
		t.Skip("test skipped for rootless containers on cgroup v1")
	}
	testContainerName := testutil.Identifier(t)[:12]
	exitedTestContainerName := fmt.Sprintf("%s-exited", testContainerName)

	base := testutil.NewBase(t)
	defer base.Cmd("rm", "-f", testContainerName).Run()
	defer base.Cmd("rm", "-f", exitedTestContainerName).Run()
	base.Cmd("run", "--name", exitedTestContainerName, testutil.AlpineImage, "echo", "'exited'").AssertOK()

	base.Cmd("run", "-d", "--name", testContainerName, testutil.AlpineImage, "sleep", "10").AssertOK()
	base.Cmd("stats", "--no-stream").AssertOutContains(testContainerName)
	base.Cmd("stats", "--no-stream", testContainerName).AssertOK()
	base.Cmd("container", "stats", "--no-stream").AssertOutContains(testContainerName)
	base.Cmd("container", "stats", "--no-stream", testContainerName).AssertOK()
}
