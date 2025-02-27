// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package count

import (
	"context"
	"fmt"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/log"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

func versionsCommand(ctx context.Context) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "Count the number of versions of specified APIs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := connection.NewClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}
			// Initialize task queue.
			taskQueue, wait := core.WorkerPool(ctx, 64)
			defer wait()
			// Generate tasks.
			name := args[0]
			if api, err := names.ParseApi(name); err == nil {
				// Iterate through a collection of APIs and count the number of versions of each.
				err = core.ListAPIs(ctx, client, api, filter, func(api *rpc.Api) {
					taskQueue <- &countVersionsTask{
						client:  client,
						apiName: api.Name,
					}
				})
				if err != nil {
					log.FromContext(ctx).WithError(err).Fatal("Failed to list APIs")
				}
			}
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filter selected resources")
	return cmd
}

type countVersionsTask struct {
	client  connection.Client
	apiName string
}

func (task *countVersionsTask) String() string {
	return "count versions " + task.apiName
}

func (task *countVersionsTask) Run(ctx context.Context) error {
	count := 0
	request := &rpc.ListApiVersionsRequest{
		Parent: task.apiName,
	}
	it := task.client.ListApiVersions(ctx, request)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		} else if err == nil {
			count++
		} else {
			return err
		}
	}
	log.Debugf(ctx, "%d\t%s", count, task.apiName)
	subject := task.apiName
	relation := "versionCount"
	artifact := &rpc.Artifact{
		Name:     subject + "/artifacts/" + relation,
		MimeType: "text/plain",
		Contents: []byte(fmt.Sprintf("%d", count)),
	}
	err := core.SetArtifact(ctx, task.client, artifact)
	if err != nil {
		return err
	}
	return nil
}
