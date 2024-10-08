module github.com/go-orb/examples/benchmarks/rps

go 1.23

toolchain go1.23.0

require (
	github.com/go-orb/go-orb v0.0.0-20240902034447-508ce5c54a46
	github.com/go-orb/plugins/client/middleware/log v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/client/orb v0.0.0-20240925070424-371b8463d2d6
	github.com/go-orb/plugins/client/orb/transport/drpc v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/grpc v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/h2c v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/hertzh2c v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/hertzhttp v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/http v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/http3 v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/client/orb/transport/https v0.0.0-20240925074005-17c0b37c3d6b
	github.com/go-orb/plugins/codecs/jsonpb v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/codecs/proto v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/codecs/yaml v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/config/source/cli/urfave v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/config/source/file v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/log/lumberjack v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/log/slog v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/registry/consul v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/registry/mdns v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/server/drpc v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/server/grpc v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/server/hertz v0.0.0-20240902051655-0791c4c590b6
	github.com/go-orb/plugins/server/http v0.0.0-20240902051655-0791c4c590b6
	github.com/google/wire v0.6.0
	github.com/hashicorp/consul/sdk v0.16.1
	google.golang.org/grpc v1.66.0
	google.golang.org/protobuf v1.34.2
	storj.io/drpc v0.0.34
)

require (
	github.com/andeya/ameda v1.5.3 // indirect
	github.com/andeya/goutil v1.0.1 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/bytedance/go-tagexpr/v2 v2.9.11 // indirect
	github.com/bytedance/gopkg v0.1.1 // indirect
	github.com/bytedance/sonic v1.12.2 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/hertz v0.9.2 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cloudwego/netpoll v0.6.3 // indirect
	github.com/cornelk/hashmap v1.0.8 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/go-orb/plugins/client/orb/transport/basehertz v0.0.0-20240925070424-371b8463d2d6 // indirect
	github.com/go-orb/plugins/client/orb/transport/basehttp v0.0.0-20240925070424-371b8463d2d6 // indirect
	github.com/go-orb/plugins/registry/regutil v0.0.0-20240902051655-0791c4c590b6 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/pprof v0.0.0-20240903155634-a8630aee4ab9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/consul/api v1.29.4 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hertz-contrib/http2 v0.1.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.62 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nyaruka/phonenumbers v1.4.0 // indirect
	github.com/onsi/ginkgo/v2 v2.20.2 // indirect
	github.com/quic-go/qpack v0.5.0 // indirect
	github.com/quic-go/quic-go v0.46.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sanity-io/litter v1.5.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/urfave/cli/v2 v2.27.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/zeebo/errs v1.3.0 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/arch v0.9.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
