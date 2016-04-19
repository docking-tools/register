#!  /bin/sh
set -e

echo " param 1 $1"
echo "param @ $@"


if [ "$1" = 'register' ]; then
    shift
    DOCKING_CONFIG=/register register \
        -ip=${HOST_IP} \
        -r=${REGISTER_URL} \
        -d=${DOCKER_URL} \
        "$@"
else
    exec "$@"
fi