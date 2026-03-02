class Ehvg < Formula
  desc "Command-line application built with Cobra"
  homepage "https://github.com/superbrave/ehvgo"
  version "2.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-darwin-arm64"
      sha256 "REPLACE_WITH_DARWIN_ARM64_SHA256"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-darwin-amd64"
      sha256 "REPLACE_WITH_DARWIN_AMD64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-linux-arm64"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-linux-amd64"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
    end
  end

  def install
    binary = if OS.mac?
      Hardware::CPU.arm? ? "ehvg-darwin-arm64" : "ehvg-darwin-amd64"
    else
      Hardware::CPU.arm? ? "ehvg-linux-arm64" : "ehvg-linux-amd64"
    end

    bin.install binary => "ehvg"
  end

  test do
    system "#{bin}/ehvg"
  end
end
