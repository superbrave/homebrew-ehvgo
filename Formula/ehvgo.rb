class Ehvgo < Formula
  desc "EHVG toolbelt"
  version "1.0.1"
  homepage "https://github.com/superbrave/ehvgo"
  license ""

  on_linux do
    on_intel do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.1-linux-amd64.gz"
      sha256 "673d577e6839a5fc1bbd2fdfb5719c8b05f63216362d73bde0019caeb3dcd5b2"
      def install
        bin.install "ehvgo-1.0.1-linux-amd64" => "ehvgo"
      end
    end

    on_arm do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.1-linux-arm64.gz"
      sha256 "b19bfaab117e6707d568fa706e6e98e5d16fcd39a984dcf59750019a5f74846e"
      def install
        bin.install "ehvgo-1.0.1-linux-arm64" => "ehvgo"
      end
    end
  end

  on_macos do
    url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-1.0.1-darwin-arm64.gz"
    sha256 "673d577e6839a5fc1bbd2fdfb5719c8b05f63216362d73bde0019caeb3dcd5b2"
    def install
      bin.install "ehvgo-1.0.1-darwin-arm64" => "ehvgo"
    end    
  end 
  test do
    system bin/"ehvgo" "version"
  end
end
