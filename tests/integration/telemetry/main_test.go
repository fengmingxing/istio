// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/ingress"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/pilot"
	"istio.io/istio/pkg/test/framework/label"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/framework/resource/environment"
)

var (
	i    istio.Instance
	p    pilot.Instance
	ingr ingress.Instance
)

func TestMain(m *testing.M) {
	framework.
		NewSuite("telemetry_test", m).
		RequireEnvironment(environment.Kube).
		RequireSingleCluster().
		Label(label.CustomSetup).
		SetupOnEnv(environment.Kube, istio.Setup(&i, func(cfg *istio.Config) {
			cfg.ControlPlaneValues = `
# Add an additional TCP port, 31400
components:
  ingressGateways:
  - name: istio-ingressgateway
    enabled: true
    k8s:
      service:
        ports:
          - port: 15020
            targetPort: 15020
            name: status-port
          - port: 80
            targetPort: 8080
            name: http2
          - port: 443
            targetPort: 8443
            name: https
          - port: 31400
            targetPort: 31400
            name: tcp
values:
  global:
    proxy:
      accessLogFile: "/dev/stdout"
  prometheus:
    enabled: true
    scrapeInterval: 5s
  grafana:
    enabled: true
  prometheus:
    enabled: true`
		})).
		Setup(func(ctx resource.Context) (err error) {
			if p, err = pilot.New(ctx, pilot.Config{}); err != nil {
				return err
			}
			if ingr, err = ingress.New(ctx, ingress.Config{
				Istio: i,
			}); err != nil {
				return err
			}
			return nil
		}).
		Run()
}
