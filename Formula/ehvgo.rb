class Ehvgo < Formula
  desc "Command-line application for managing EHVG resources"
  homepage "https://github.com/superbrave/homebrew-ehvgo"
  version "2.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-darwin-arm64"
      sha256 "65d41a67e1f3eb249f648fb4c5b7e31defe005a19fd20d7f2097a4724f38ee16"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-darwin-amd64"
      sha256 "c4adc96625c2646d0c9b4ee0737e20d6a0f35725cca77f9a36ae81635c520956"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-linux-arm64"
      sha256 "cf63b65efad9216bd4ab01e88fce1e15811e19111deaa25889cfec169d7ba05e"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvgo-linux-amd64"
      sha256 "e7ecc013ee767e2d6d2b99c6852f761826841ca00cfaa04740454d9aac741279"
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
