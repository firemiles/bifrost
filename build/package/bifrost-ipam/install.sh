#!/bin/sh
# https://github.com/projectcalico/cni-plugin/blob/e1a1ea6834cad8235db6d0e87901cd144d60e265/k8s-install/scripts/install-cni.sh

SERVICE_ACCOUNT_PATH=/var/run/secrets/kubernetes.io/serviceaccount
KUBE_CA_FILE=${KUBE_CA_FILE:-$SERVICE_ACCOUNT_PATH/ca.crt}
SKIP_TLS_VERIFY=${SKIP_TLS_VERIFY:-false}
# Pull out service account token.
SERVICEACCOUNT_TOKEN=$(cat $SERVICE_ACCOUNT_PATH/token)
# TLS_CFG needs to be set regardless (to avoid error when creating calico-kubeconfig)
TLS_CFG=""

# Check if we're running as a k8s pod.
if [ -f "$SERVICE_ACCOUNT_PATH/token" ]; then
  # We're running as a k8d pod - expect some variables.
  if [ -z "${KUBERNETES_SERVICE_HOST}" ]; then
    echo "KUBERNETES_SERVICE_HOST not set"; exit 1;
  fi
  if [ -z "${KUBERNETES_SERVICE_PORT}" ]; then
    echo "KUBERNETES_SERVICE_PORT not set"; exit 1;
  fi

  if [ "$SKIP_TLS_VERIFY" = "true" ]; then
    TLS_CFG="insecure-skip-tls-verify: true"
  elif [ -f "$KUBE_CA_FILE" ]; then
    TLS_CFG="certificate-authority-data: $(base64 -w 0 "$KUBE_CA_FILE")"
  fi

  # Write a kubeconfig file for the CNI plugin.  Do this
  # to skip TLS verification for now.  We should eventually support
  # writing more complete kubeconfig files. This is only used
  # if the provided CNI network config references it.
  touch /host/etc/cni/net.d/bifrost-kubeconfig
  chmod "${KUBECONFIG_MODE:-600}" /host/etc/cni/net.d/bifrost-kubeconfig
  cat > /host/etc/cni/net.d/bifrost-kubeconfig <<EOF
# Kubeconfig file for Bifrost CNI plugin.
apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: ${KUBERNETES_SERVICE_PROTOCOL:-https}://[${KUBERNETES_SERVICE_HOST}]:${KUBERNETES_SERVICE_PORT}
    $TLS_CFG
users:
- name: bifrost
  user:
    token: "${SERVICEACCOUNT_TOKEN}"
contexts:
- name: bifrost-context
  context:
    cluster: local
    user: bifrost
current-context: bifrost-context
EOF

fi

cp /bifrost-ipam /opt/cni/bin
chmod +x /opt/cni/bin/bifrost-ipam