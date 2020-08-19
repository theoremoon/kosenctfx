#!/bin/sh
vagrant up
vagrant ssh-config > ssh_config
cat <<EOF > ansible.cfg
[default]
hostfile = hosts

[ssh_connection]
ssh_args = -F ssh_config
EOF

cat <<EOF > hosts
[all]
host1
host2
host3
host4
EOF

ansible-playbook -i hosts ansible/docker.yml
