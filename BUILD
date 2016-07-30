load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_binary", "go_library")

package(
    default_visibility = ["//:__subpackages__"],
)

go_prefix("github.com/cceckman/primes")

go_library(
    name = "go_default_library",
    srcs = [
        "erat.go",
        "interface.go",
    ],
)

go_binary(
    name = "benchmark",
    srcs = ["main.go"],
    deps = [":go_default_library"],
)
