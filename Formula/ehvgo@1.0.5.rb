class Ehvgo105 < Formula
  desc "EHVG toolbelt"
  version "1.0.5"
  homepage "https://github.com/superbrave/ehvgo"
  license ""

  on_linux do
    on_intel do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.5-linux-amd64.gz"
      sha256 "dfedd03d2e48861bbfaeef6fc2d7295481c516fb1c3bbca050f1db4d13a5a4e7"
      def install
        bin.install "ehvgo-1.0.5-linux-amd64" => "ehvgo"
      end
    end

    on_arm do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.5-linux-arm64.gz"
      sha256 "64e86fab0c6c9aab50b6a635fafb2fdb9605712a618eb5e79cb5798e54043952"
      def install
        bin.install "ehvgo-1.0.5-linux-arm64" => "ehvgo"
      end
    end
  end

  on_macos do
    url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.5-darwin-arm64.gz"
    sha256 "1d96a35e1eab3f0d8879297c7502429fdef8d569dfcbf95aa991ff1bac1e2839"
    def install
      bin.install "ehvgo-1.0.5-darwin-arm64" => "ehvgo"
    end    
  end 
  test do
    system bin/"ehvgo" "version"
  end
end
