load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "events.go",
        "scheduling_queue.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/scheduler/internal/queue",
    visibility = ["//pkg/scheduler:__subpackages__"],
    deps = [
        "//pkg/scheduler/framework:go_default_library",
        "//pkg/scheduler/internal/heap:go_default_library",
        "//pkg/scheduler/metrics:go_default_library",
        "//pkg/scheduler/util:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/types:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//staging/src/k8s.io/client-go/listers/core/v1:go_default_library",
        "//staging/src/k8s.io/client-go/tools/cache:go_default_library",
        "//vendor/k8s.io/klog/v2:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["scheduling_queue_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//pkg/api/v1/pod:go_default_library",
        "//pkg/scheduler/framework:go_default_library",
        "//pkg/scheduler/framework/plugins/queuesort:go_default_library",
        "//pkg/scheduler/metrics:go_default_library",
        "//pkg/scheduler/util:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/types:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/clock:go_default_library",
        "//staging/src/k8s.io/client-go/informers:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//staging/src/k8s.io/component-base/metrics/testutil:go_default_library",
        "//vendor/github.com/google/go-cmp/cmp:go_default_library",
        "//vendor/github.com/google/go-cmp/cmp/cmpopts:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
