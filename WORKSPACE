workspace(name = "com_github_Xjs_gopher_storer")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "7949fc6cc17b5b191103e97481cf8889217263acf52e00b560683413af204fcb",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.16.0/bazel-gazelle-0.16.0.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_go_sql_driver_mysql",
    importpath = "github.com/go-sql-driver/mysql",
    sha256 = "ebf7924f6aa009a708b218015f99c9466e5a0c8e1b24f6e9bae9e1cd7a505710",
    strip_prefix = "mysql-1.4.0",
    urls = ["https://github.com/go-sql-driver/mysql/archive/v1.4.0.zip"],
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sha256 = "6f9b8cdf96725fad0fe750c8aad10105c91805e0c76931dccbdec3b7f6b1bcbf",
    strip_prefix = "mux-1.6.2",
    urls = ["https://github.com/gorilla/mux/archive/v1.6.2.zip"],
)

go_repository(
    name = "com_github_gorilla_context",
    commit = "51ce91d2eaddeca0ef29a71d766bb3634dadf729",
    importpath = "github.com/gorilla/context",
)
