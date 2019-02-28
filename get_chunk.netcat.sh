#!/bin/bash

usage() { echo "Usage: $0 [-f <string>] [-i <int>] [-s <int>] [-p <int>]" 1>&2; exit 1; }

file=""
index=0
size=1024
port=3333


while getopts ":f:i:s:p:" o; do
    case "${o}" in
        f)
            file=${OPTARG}
            ;;
        i)
            index=${OPTARG}
            ;;
        s)
            size=${OPTAG}
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

if [ -z "${file}" ] || [ -z "${index}" ] || [ -z "${size}" ] || [ -z "${port}" ]; then
    usage
fi

echo "{\"Code\":4,\"Info\":{\"FileName\":\"$file\",\"ChunkIndex\":$index,\"ChunkSize\":$size}}" | nc localhost $port
