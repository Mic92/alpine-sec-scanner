# alpine-sec-scanner

Proof-of-concept application for [vmsh](https://github.com/Mic92/vmsh)

Checks installed packages against https://secdb.alpinelinux.org/

## Usage

```console
$ go build
# scan the systems package for insecure packages
$ ./alpine-sec-scanner /
# test alpine installation found in /some-chroot
$ ./alpine-sec-scanner /some-chroot
```
