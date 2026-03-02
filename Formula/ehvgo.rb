class Ehvg < Formula
  desc "Command-line application for managing EHVG resources"
  homepage "https://github.com/superbrave/homebrew-ehvgo"
  version "2.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-darwin-arm64"
      sha256 "5bea7b6f422766f3508534662a0bbefd254bd31ba08a61426e790d73a85b1fd7"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-darwin-amd64"
      sha256 "70bdec4d22cac2f6438c5819ca1cf965fefcf2b4528200c8652428c4f47a764d"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-linux-arm64"
      sha256 "33bb43f61c21754836507bdf52fe34b680cc56f9fbf6875a9175d424039d791e"
    else
      url "https://github.com/superbrave/ehvgo/releases/download/#{version}/ehvg-linux-amd64"
      sha256 "5ada982e8d7274c842a51b4ad28ed4a4c5fa1a56958116dd9a61a8852cd864a4"
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
