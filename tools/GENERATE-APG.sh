#!/bin/bash
#
# Copyright 2021 Google LLC. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

source tools/PROTOS.sh
clone_common_protos

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_cli@latest

echo "Generating Go client CLI for ${SERVICE_PROTOS[@]}"
protoc ${SERVICE_PROTOS[*]} \
	--proto_path='.' \
	--proto_path=$COMMON_PROTOS_PATH \
  	--go_cli_opt='root=apg' \
  	--go_cli_opt='gapic=github.com/apigee/registry/gapic' \
  	--go_cli_out='cmd/apg'
