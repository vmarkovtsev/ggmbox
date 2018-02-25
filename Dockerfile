FROM ubuntu:16.04

ADD ggmbox.py /
ADD parse.go /root/go/src/github.com/vmarkovtsev/ggmbox/parse.go

RUN apt-get update && \
    apt-get install -y ca-certificates locales python3 wget python3-dev gcc git && \
    wget -O - https://bootstrap.pypa.io/get-pip.py | python3 && \
    pip3 install --no-cache-dir scrapy && \
    pip3 uninstall -y pip && \
    wget -O - https://dl.google.com/go/go1.10.linux-amd64.tar.gz | tar -C /usr/local -xz && \
    export PATH=$PATH:/usr/local/go/bin && \
    mkdir /root/go/src/github.com/vmarkovtsev/ggmbox/.git && \
    go get github.com/vmarkovtsev/ggmbox && \
    rm -rf /usr/local/go && \
    mv /root/go/bin/ggmbox /usr/local/bin/parse && \
    apt-get remove -y wget python3-dev gcc git && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /root/go && \
    locale-gen en_US.UTF-8
