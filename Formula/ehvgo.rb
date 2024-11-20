class Ehvgo < Formula
  desc "EHVG toolbelt"
  version "0.1.13"
  homepage "https://github.com/superbrave/ehvgo"
  url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-0.1.13-linux-amd64.gz"
  sha256 "02a79dd745c1d658b5c8a373845450afa1f5b58734d3f9c5d8f4f12919005f07"
  license ""

  on_linux do
    on_intel do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-0.1.13-linux-amd64.gz"
      sha256 "02a79dd745c1d658b5c8a373845450afa1f5b58734d3f9c5d8f4f12919005f07"
      def install
        bin.install "ehvgo-0.1.13-linux-amd64" => "ehvgo"
      end
    end
    on_arm do
      url "https://dist.ehealthsystems.nl/ehvgo/ehvgo-0.1.13-linux-arm64.gz"
      sha256 "e6b1f8665d8924f52c5603cf7114c2db489dfc7f048458140ebe090ed177c529"
      def install
        bin.install "ehvgo-0.1.13-linux-arm64" => "ehvgo"
      end
    end
  end

  test do
    # `test do` will create, run in and delete a temporary directory.
    #
    # This test will fail and we won't accept that! For Homebrew/homebrew-core
    # this will need to be a test that verifies the functionality of the
    # software. Run the test with `brew test ehvgo`. Options passed
    # to `brew install` such as `--HEAD` also need to be provided to `brew test`.
    #
    # The installed folder is not in the path, so use the entire path to any
    # executables being tested: `system bin/"program", "do", "something"`.
    system bin/"ehvgo" "version"
  end
end
