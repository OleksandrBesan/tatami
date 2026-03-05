class Tatami < Formula
  desc "Terminal workspace manager with Zellij/Tmux integration"
  homepage "https://github.com/OleksandrBesan/tatami"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/OleksandrBesan/tatami/releases/download/v#{version}/tatami_#{version}_darwin_arm64.tar.gz"
      # sha256 "PLACEHOLDER" # Update after release
    else
      url "https://github.com/OleksandrBesan/tatami/releases/download/v#{version}/tatami_#{version}_darwin_amd64.tar.gz"
      # sha256 "PLACEHOLDER" # Update after release
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/OleksandrBesan/tatami/releases/download/v#{version}/tatami_#{version}_linux_arm64.tar.gz"
      # sha256 "PLACEHOLDER" # Update after release
    else
      url "https://github.com/OleksandrBesan/tatami/releases/download/v#{version}/tatami_#{version}_linux_amd64.tar.gz"
      # sha256 "PLACEHOLDER" # Update after release
    end
  end

  def install
    bin.install "tatami"
  end

  def caveats
    <<~EOS
      To enable automatic cd, add to your ~/.zshrc or ~/.bashrc:

        tatami() {
          local output
          output=$(TATAMI_WRAPPER=1 command tatami "$@")
          local exit_code=$?
          if [[ $exit_code -eq 0 && -d "$output" ]]; then
            cd "$output"
          elif [[ -n "$output" ]]; then
            echo "$output"
          fi
          return $exit_code
        }
    EOS
  end

  test do
    system "#{bin}/tatami", "--version"
  end
end
