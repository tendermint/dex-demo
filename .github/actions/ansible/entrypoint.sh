#!/bin/sh

mkdir -p ~/.ssh && echo $INPUT_SSH_KEY_B64 | base64 -d > ~/.ssh/id_rsa && chmod 600 ~/.ssh/id_rsa

if [ "$INPUT_INVENTORY_SCRIPT" == "true" ]
then
    chmod +x $INPUT_INVENTORY;
fi

if [ -n "$INPUT_VAULT_PASSWORD" ]; then
    echo $INPUT_VAULT_PASSWORD > /tmp/vault-passwd
    export ANSIBLE_VAULT_PASSWORD_FILE=/tmp/vault-passwd
fi

sleep $INPUT_WAIT
ANSIBLE_HOST_KEY_CHECKING=False \
    ansible-playbook \
    -i $INPUT_INVENTORY \
    $INPUT_PLAYBOOK \
    -e "$INPUT_EXTRA_VARS" \
    --key-file ~/.ssh/id_rsa
