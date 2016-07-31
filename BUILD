load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library", "go_test")

package(
    default_visibility = ["//:__subpackages__"],
)

go_prefix("github.com/cceckman/primes")

go_library(
    name = "go_default_library",
    srcs = [
        "erat.go",
        "primes.go",
        "parrerat.go",
    ],
)


#go_test(
#    name = "primes_test",
#    srcs = ["primes_test.go"],
#    library = ":go_default_library",
#)
