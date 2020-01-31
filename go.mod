module github.com/firemiles/bifrost

go 1.13

require (
	github.com/containernetworking/cni v0.7.1
	github.com/containernetworking/plugins v0.8.5
	github.com/go-logr/logr v0.1.0
	github.com/golang/protobuf v1.3.2
	github.com/j-keck/arping v1.0.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	golang.org/x/sys v0.0.0-20190626221950-04f50cda93cb // indirect
	google.golang.org/grpc v1.26.0
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.49.0
	cloud.google.com/go/bigquery => github.com/googleapis/google-cloud-go/bigquery v1.3.0
	cloud.google.com/go/datastore => github.com/googleapis/google-cloud-go/datastore v1.0.0
	cloud.google.com/go/pubsub => github.com/googleapis/google-cloud-go/pubsub v1.1.0
	cloud.google.com/go/storage => github.com/googleapis/google-cloud-go/storage v1.4.0
	golang.org/x/arch => github.com/golang/arch v0.0.0-20191126211547-368ea8f32fff
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191128160524-b544559bb6d1
	golang.org/x/exp => github.com/golang/exp v0.0.0-20191129062945-2f5052295587
	golang.org/x/image => github.com/golang/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/lint => github.com/golang/lint v0.0.0-20191125180803-fdd1cda4f05f
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20191123054942-d9e324ca8c38
	golang.org/x/mod => github.com/golang/mod v0.1.0
	golang.org/x/net => github.com/golang/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20191122200657-5d9234df094c
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20191128015809-6d18c012aee9
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191130070609-6e064ea0cf2d
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.14.0
	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20191223191004-3caeed10a8bf
	google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0
	gopkg.in/inf.v0 => github.com/go-inf/inf v0.9.1
	gopkg.in/mgo.v2 => github.com/go-mgo/mgo v0.0.0-20180705113738-7446a0344b78
	gopkg.in/yaml.v2 v2.2.4 => github.com/go-yaml/yaml v0.0.0-20191002183336-f221b8435cfb
	gopkg.in/yaml.v3 => github.com/go-yaml/yaml v0.0.0-20191120175047-4206685974f2
	k8s.io/api => github.com/kubernetes/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apimachinery => github.com/kubernetes/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go => github.com/kubernetes/client-go v0.0.0-20191114101535-6c5935290e33
	k8s.io/code-generator => github.com/kubernetes/code-generator v0.0.0-20191121015212-c4c8f8345c7e
	k8s.io/klog v0.4.0 => github.com/kubernetes/klog v0.4.0
	k8s.io/kube-openapi => github.com/kubernetes/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/utils => github.com/kubernetes/utils v0.0.0-20200124190032-861946025e34
	sigs.k8s.io/controller-runtime v0.4.0 => github.com/kubernetes-sigs/controller-runtime v0.4.0
	sigs.k8s.io/controller-tools/cmd/controller-gen v0.2.4 => github.com/kubernetes-sigs/controller-tools/cmd/controller-gen v0.2.4
	sigs.k8s.io/testing_frameworks => github.com/kubernetes-sigs/testing_frameworks v0.1.2
	sigs.k8s.io/yaml v1.1.0 => github.com/kubernetes-sigs/yaml v1.1.0
)
