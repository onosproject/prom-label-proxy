module github.com/prometheus-community/prom-label-proxy

go 1.16

require (
	github.com/efficientgo/tools/core v0.0.0-20210201224146-3d78f4d30648
	github.com/go-openapi/runtime v0.19.29
	github.com/go-openapi/strfmt v0.20.1
	github.com/google/gnxi v0.0.0-20210716134716-cb5c55758a07
	github.com/onosproject/config-models/modelplugin/plproxy-1.0.0 v0.8.1
	github.com/onosproject/onos-lib-go v0.7.20
	github.com/onosproject/sdcore-adapter v0.1.36
	github.com/openconfig/gnmi v0.0.0-20210707145734-c69a5df04b53
	github.com/openconfig/ygot v0.12.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/alertmanager v0.23.0
	github.com/prometheus/prometheus v1.8.2-0.20210811141203-dcb07e8eac34
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.39.0
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/onosproject/config-models/modelplugin/plproxy-1.0.0 => ../config-models/modelplugin/plproxy-1.0.0
	github.com/onosproject/sdcore-adapter => ../sdcore-adapter

)
