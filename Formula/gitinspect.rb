class Gitinspect < Formula
  desc "Turn any Git repo into an AI-friendly, token-efficient snapshot"
  homepage "https://github.com/richie-rich90454/gitinspect"
  url "https://github.com/richie-rich90454/gitinspect/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "Apache-2.0"
  head "https://github.com/richie-rich90454/gitinspect.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/gitinspect"
  end

  test do
    system bin/"gitinspect", "--help"
  end
end
