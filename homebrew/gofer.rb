require 'formula'

class Gofer < Formula
  VERSION = "0.0.3"

  homepage "https://github.com/chuckpreslar/gofer"
  head "https://github.com/chuckpreslar/gofer.git"
  url "https://github.com/chuckpreslar/gofer/releases/download/v#{VERSION}/gofer_#{VERSION}.tar.gz"
  version VERSION

  def install
    unless system "echo $GOPATH"
      abort "To use gofer, you must first set your $GOPATH environment variable set."     
    end

    unless system "go get -u github.com/chuckpreslar/gofer"
      abort "Failed to install gofer package."
    end

    bin.install "gofer"
  end

end
