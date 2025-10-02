beanstool [![test](https://github.com/rafaelespinoza/beanstool/actions/workflows/test.yaml/badge.svg)](https://github.com/rafaelespinoza/beanstool/actions/workflows/test.yaml)
==============================

Dependency free [beanstalkd](https://beanstalkd.github.io/) admin tool.

This is a maintained fork of https://github.com/src-d/beanstool/. That is a rework of the wonderful [bsTools](https://github.com/jimbojsb/bstools) with some extra features and of course without need to install any dependency. Very useful in companion of the server in a small docker container.

Installation
------------

If you have the [`gh` CLI](https://cli.github.com/) installed:
```sh
# Target this fork with the flag, `-R | --repo`
$ gh --repo rafaelespinoza/beanstool release list

# Pick the build for your platform. In this example, it's linux amd64.
$ gh --repo rafaelespinoza/beanstool release download -p '*linux_amd64*' -p checksums.txt -D /tmp

# [Optional] Verify the checksum downloaded in previous step from download destination dir
$ cd /tmp
$ grep linux_amd64 checksums.txt | sha256sum -c

# Unarchive and copy to your bin place. The ${version} below is a placeholder.
$ tar -xvzf beanstool_${version}_linux_amd64.tar.gz
$ cp beanstool_${version}_linux_amd64/beanstool /usr/local/bin/
```

Browse the [`releases`](https://github.com/rafaelespinoza/beanstool/releases) section to see other archs and versions.

Usage
-----

```sh
Usage:
  beanstool [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  bury   bury existing jobs from ready state
  kick   kicks jobs from buried back into ready
  delete a job from a queue
  peek   peeks a job from a queue
  put    put a job into a tube
  stats  print stats on all tubes
  tail   tails a tube and prints his content
```

As example this is the output of the command `./beanstool stats`:

```
+---------+----------+----------+----------+----------+----------+---------+-------+
| Name    | Buried   | Delayed  | Ready    | Reserved | Urgent   | Waiting | Total |
+---------+----------+----------+----------+----------+----------+---------+-------+
| foo     | 20       | 0        | 5        | 0        | 0        | 0       | 28    |
| default | 0        | 0        | 0        | 0        | 0        | 0       | 0     |
+---------+----------+----------+----------+----------+----------+---------+-------+
```

License
-------

MIT, see [LICENSE](LICENSE)
