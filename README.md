ggmbox [![Build Status](https://travis-ci.org/vmarkovtsev/ggmbox.svg?branch=master)](https://travis-ci.org/vmarkovtsev/ggmbox) [![Build status](https://ci.appveyor.com/api/projects/status/x57poug9apd0bs2h?svg=true)](https://ci.appveyor.com/project/vmarkovtsev/ggmbox) [![Docker Build Status](https://img.shields.io/docker/build/vmarkovtsev/ggmbox.svg)](https://hub.docker.com/r/vmarkovtsev/ggmbox)
======

Google Groups raw emails crawler and parser. Turbo speed and reliable!
The downloaded messages are in [RFC 822](https://www.ietf.org/rfc/rfc822.txt) format - taken verbatim
from the Google servers. 

### Installation

#### Docker

Docker is the simplest option. Go to [![DockerHub](https://img.shields.io/docker/build/vmarkovtsev/ggmbox.svg)](https://hub.docker.com/r/vmarkovtsev/ggmbox)
Prepend `docker run -it --rm vmarkovtsev/ggmbox` to all the commands in the "Usage" section.

#### Crawler

Requirements: [Python 3](https://www.python.org/) and [Scrapy](https://scrapy.org/). Download
[`ggmbox.py`](ggmbox.py) file.

#### Parser

Requirements: [Go](https://golang.org/).

```
go get -v github.com/vmarkovtsev/ggmbox
```

### Usage

#### Crawler

```
scrapy runspider -a name=golang-nuts -o result.json -t json ggmbox.py
```

Replace "golang-nuts" with the actual group name. The raw emails will be saved by default to the
corresponding directory.

#### Parser

```
./parse golang-nuts > dataset.csv
```

Replace "golang-nuts" with the actual directory name with raw emails. The plain text threads will
be written to `dataset.csv`, one thread per line. Special characters are escaped.

### Performance

#### Crawler

[golang-nuts](https://groups.google.com/d/forum/golang-nuts) group was fully fetched on 24/02/2018 with
30043 topics and 192654 messages **in 3 hours** at 1gbps connection speed.
The raw emails occupied 1.6 GB on disk.

Compare to 1 day using [icy/google-group-crawler](https://github.com/icy/google-group-crawler),
it fetched only 63% and then stopped without any errors reported, or to
[henryk/gggd](https://github.com/henryk/gggd), it fetched only 3% within one hour and then
unexpectedly stopped too.

#### Parser

It takes **7 seconds** to parse 1.6 GB of raw emails on a 32-core machine.

### Contributions

...are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

### License

[MIT](LICENSE).
