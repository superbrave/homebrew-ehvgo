class Ehvgo < Formula
  desc "Command-line application for managing EHVG resources"
  homepage "https://github.com/superbrave/homebrew-ehvgo"
  version "2.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-darwin-arm64"
      sha256 "f612b60899b83fd619d9f94f506b01e68a0fb134263ea1701c8a33a722a46045"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-darwin-amd64"
      sha256 "ba94d41accbe1815930a183a68ec7e6a39cb22a0662dac1e36438de64b7655fe"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-linux-arm64"
      sha256 "132588f4ed658b3c8d0bcddbcaf0a4f670d3b30bb219b8738dc4ae87c7fd8991"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-linux-amd64"
      sha256 "1ffa87664de1698ef1fb99b5ef36b9ec4257760372801bfe5291c6abc3604a09"
    end
  end

  def install
    binary = if OS.mac?
      Hardware::CPU.arm? ? "ehvgo-darwin-arm64" : "ehvgo-darwin-amd64"
    else
      Hardware::CPU.arm? ? "ehvgo-linux-arm64" : "ehvgo-linux-amd64"
    end

    bin.install binary => "ehvgo"
  end

  test do
    system "#{bin}/ehvgo"
  end
end
