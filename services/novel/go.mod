module novel

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/codahale/hdrhistogram v0.0.0-00010101000000-000000000000 // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0
	github.com/golang/protobuf v1.4.2
	github.com/ilylx/gconv v0.0.0-20240713143307-cc305f890dcd
	github.com/micro/go-micro/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/scg130/tools v0.10.9
	github.com/sirupsen/logrus v1.4.2
	google.golang.org/protobuf v1.25.0 // indirect
	xorm.io/xorm v1.3.9
)

replace github.com/codahale/hdrhistogram => github.com/HdrHistogram/hdrhistogram-go v0.9.0
