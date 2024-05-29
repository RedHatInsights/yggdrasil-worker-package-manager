#!/bin/sh

set -x

# Dump shell environment for context
env

# Print relevant package info
rpm -qi yggdrasil
rpm -qi yggdrasil-worker-package-manager

# Configure yggdrasil for local-dispatch only
sed -i -e 's/protocol = .*$/protocol = "none"/' \
       -e '/server = .*/d' \
       -e 's/log-level = .*/log-level = "debug"/' \
       /etc/yggdrasil/config.toml
echo 'message-journal = ":memory:"' >> /etc/yggdrasil/config.toml
sed -i -e 's/log-level = .*/log-level = "debug"/' /etc/yggdrasil-worker-package-manager/config.toml

# Ensure yggdrasil is running
systemctl start yggdrasil
systemctl status yggdrasil
busctl --system status com.redhat.Yggdrasil1

# Locally dispatch a command to install the 'vim' package, triggering the
# bus-activated worker service
echo '{"command":"install","name":"vim"}' | yggctl dispatch --worker package_manager -

# Verify the worker started via bus-activation and can be introspected
WORKER_UNIT=$(busctl --system status com.redhat.Yggdrasil1.Worker1.package_manager | grep ^Unit= | cut -f2 -d=)
busctl --system introspect com.redhat.Yggdrasil1.Worker1.package_manager /com/redhat/Yggdrasil1/Worker1/package_manager com.redhat.Yggdrasil1.Worker1
yggctl workers list
