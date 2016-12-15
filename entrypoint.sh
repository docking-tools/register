#!  /bin/sh
set -e

echo " param 1 $1"
echo "param @ $@"


if [ "$1" = 'register' ]; then
    shift
    DOCKING_CONFIG=/ register \
        -ip=${HOST_IP} \
        -d=${DOCKER_URL} \
        "$@"
else
    exec "$@"
fi