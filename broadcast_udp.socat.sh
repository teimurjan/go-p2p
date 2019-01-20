#!/bin/bash

usage() { echo "Usage: $0 [-m <string>] [-p <int>]" 1>&2; exit 1; }

port=3333
message=""


while getopts ":m:p:" o; do
    case "${o}" in
        m)
            message=${OPTARG}
            ;;
        p)
            port=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ -z "${message}" ] || [ -z "${port}" ]; then
    usage
fi

echo $message | socat - UDP4-DATAGRAM:255.255.255.255:$port,broadcast