FROM container-registry.oracle.com/database/instantclient:12.2.0.1

MAINTAINER Chris Duncan <chris.duncan@returnpath.com>

# Oracle environment variables, pkg config and tns settings
ENV ORACLE_HOME /usr/lib/oracle/12.2/client64
ENV ORACLE_INCLUDE /usr/include/oracle/12.2/client64
ENV LD_LIBRARY_PATH ${ORACLE_HOME}/lib
ENV PKG_CONFIG_PATH ${ORACLE_INCLUDE}:/usr/lib64/pkgconfig:/usr/share/pkgconfig
ENV TNS_ADMIN ${ORACLE_HOME}/network/admin
ENV NLS_LANG American_America.AL32UTF8
ENV PATH ${ORACLE_HOME}/bin:${PATH}

# TODO: maybe delete this?
RUN set -eux; \
  mkdir -p "${TNS_ADMIN}" "${LD_LIBRARY_PATH}/sdk"; \
  ln -s "${ORACLE_INCLUDE}" "${LD_LIBRARY_PATH}/sdk/include"; \
  ln -s "${ORACLE_HOME}/lib" "${LD_LIBRARY_PATH}/sdk/lib";

## Above is needed for production. Below is needed for Dev ##

# Install requirements for Golang
RUN set -eux; \
  yum -y update; \
  yum -y install curl tar gzip gcc-c++ gcc glibc-devel make pkgconfig; \
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

# TODO MOVE UP ABOVE
COPY ./docker/golang_oracle_client/oci8.pc ${ORACLE_INCLUDE}/
COPY ./docker/golang_oracle_client/*.ora ${TNS_ADMIN}/

CMD ["go"]
