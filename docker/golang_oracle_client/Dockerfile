FROM container-registry.oracle.com/database/instantclient:12.2.0.1

MAINTAINER Chris Duncan <chris.duncan@returnpath.com>

# Oracle environment variables and pkg-config
ENV ORACLE_HOME /usr/lib/oracle/12.2/client64
ENV ORACLE_INCLUDE /usr/include/oracle/12.2/client64
ENV LD_LIBRARY_PATH ${ORACLE_HOME}/lib
ENV PKG_CONFIG_PATH ${ORACLE_INCLUDE}:/usr/lib64/pkgconfig:/usr/share/pkgconfig
ENV NLS_LANG American_America.AL32UTF8
ENV PATH ${ORACLE_HOME}/bin:${PATH}

COPY ./docker/golang_oracle_client/oci8.pc ${ORACLE_INCLUDE}/

## Above is needed for Production (until we figure out how to statically compile) ##
## Below is needed for Development and Building only ##

# Install requirements for Golang
RUN set -eux; \
  yum-config-manager --enable ol7_optional_latest; \
  yum -y update; \
  yum -y install curl tar gzip gcc-c++ gcc glibc-devel make pkgconfig glibc-static; \
  yum clean all;

# Install Golang
ENV GO_VERSION 1.8.3
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:${GOROOT}/bin:${PATH}

RUN set -eux; \
  mkdir -p ${GOROOT} ${GOPATH}/src ${GOPATH}/bin; \
  curl -SLO "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"; \
  echo "1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772 go${GO_VERSION}.linux-amd64.tar.gz" | sha256sum -c -; \
  tar -xzf "go${GO_VERSION}.linux-amd64.tar.gz" -C ${GOROOT} --strip-components=1; \
  chmod -R 777 "${GOPATH}"; \
  rm -f "go${GO_VERSION}.linux-amd64.tar.gz"; \
  go version;

CMD ["go"]
