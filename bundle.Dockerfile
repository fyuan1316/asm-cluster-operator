FROM scratch

LABEL operators.operatorframework.io.bundle.mediatype.v1=plain
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=test-operator
LABEL operators.operatorframework.io.bundle.channels.v1=stable,beta
LABEL operators.operatorframework.io.bundle.channel.default.v1=stable

COPY bundle /manifests/
COPY metadata /metadata/
