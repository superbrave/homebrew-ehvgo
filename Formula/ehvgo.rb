class Ehvgo < Formula
  desc "EHVG toolbelt"
  version "1.0.0"
  homepage "https://github.com/superbrave/ehvgo"
  license ""

  on_linux do
    on_intel do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.0-linux-amd64.gz"
      sha256 "c10cb6447b813f5947a603c1c82d7eb4b67c3ae9143d33316111225d915db85e"
      def install
        bin.install "ehvgo-1.0.0-linux-amd64" => "ehvgo"
      end
    end

    on_arm do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.0-linux-arm64.gz"
      sha256 "34d99d94066d2b09c7e4f82078e82970b155f5558113cf0fd0748d9f5cfdea4a"
      def install
        bin.install "ehvgo-1.0.0-linux-arm64" => "ehvgo"
      end
    end
  end

  on_macos do
    url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.0-darwin-arm64.gz"
    sha256 "970d0578cd1a392a8ed5b9db2456b6a0fc899134ac69855a6cebe85050b2406a"
    def install
      bin.install "ehvgo-1.0.0-darwin-arm64" => "ehvgo"
    end    
  end 
  test do
    system bin/"ehvgo" "version"
  end
end
