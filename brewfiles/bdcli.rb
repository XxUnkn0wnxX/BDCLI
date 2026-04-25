class Bdcli < Formula
  desc "BetterDiscord CLI fork for macOS 11 Big Sur on older Intel Macs"
  homepage "https://github.com/XxUnkn0wnxX/BDCLI"
  license "Apache-2.0"
  head "https://github.com/XxUnkn0wnxX/BDCLI.git", branch: "main"

  depends_on "xxunkn0wnxx/bdcli/go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ENV["GOOS"] = "darwin"
    ENV["GOARCH"] = "amd64"
    ENV["GOAMD64"] = "v1"
    ENV["MACOSX_DEPLOYMENT_TARGET"] = "11.0"

    short_sha = Utils.safe_popen_read("git", "rev-parse", "--short", "HEAD").chomp
    build_date = Time.now.utc.strftime("%Y-%m-%dT%H:%M:%SZ")
    ldflags = %W[
      -s -w
      -X main.version=#{short_sha}
      -X main.commit=#{short_sha}
      -X main.date=#{build_date}
    ]

    system "go", "build", *std_go_args(output: bin/"bdcli", ldflags: ldflags), "./main.go"
    generate_completions_from_executable(bin/"bdcli", "completion")

    pkgshare.install "README.md", "LICENSE"
  end

  test do
    assert_match "bdcli", shell_output("#{bin}/bdcli --help")
    assert_match "Print the version number", shell_output("#{bin}/bdcli help version")
  end
end
