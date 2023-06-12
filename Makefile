#----------------------------------------------------------------------------
#  SPDX-License-Identifier: Apache-2.0
#
#  Copyright (c) 2023 Cisco Systems, Inc.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#----------------------------------------------------------------------------


.PHONY: license-check
license-check:
	docker run -it -v ${PWD}:/src -w /src ghcr.io/google/addlicense -check -l apache -c "Cisco" -ignore '**/Dockerfile' -ignore '**/*.yaml' -ignore '**/*.xml' -ignore '**/node_modules/**' -ignore 'src/frontend/.*/**' -ignore 'src/frontend/src/quasar.d.ts' .

.PHONY: license-add-spdx
license-add-spdx:
	docker run -it -v ${PWD}:/src -w /src ghcr.io/google/addlicense -check -l apache -c "Cisco" -ignore '**/Dockerfile' -ignore '**/*.yaml' -ignore '**/*.xml' -ignore '**/node_modules/**' -ignore 'src/frontend/.*/**' -ignore 'src/frontend/src/quasar.d.ts' -s .

.PHONY: scan
scan:
	~/Downloads/kubeclarity-cli-2.18.1-darwin-amd64/kubeclarity-cli analyze src/frontend/node_modules --input-type dir -o frontend-scan.sbom
	~/Downloads/kubeclarity-cli-2.18.1-darwin-amd64/kubeclarity-cli scan frontend-scan.sbom --input-type sbom
