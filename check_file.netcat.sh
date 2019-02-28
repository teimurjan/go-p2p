#!/bin/bash

usage() { echo "Usage: $0 [-f <string>] [-p <int>]" 1>&2; exit 1; }

file=""
port=3333


while getopts ":f:p:" o; do
    case "${o}" in
        f)
            file=${OPTARG}
            ;;
        p)
            port=${OPTAG}
            ;;
        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ -z "${file}" ] || [ -z "${port}" ]; then
    usage
fi

echo "{\"Code\":3,\"Info\":{\"FileName\":\"$file\"}}" | nc localhost $port