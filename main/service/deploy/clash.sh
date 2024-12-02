#!/bin/env bash


BIN_PATH=$(dirname "$0")

CONFIG_FILE=${BIN_PATH}/config/token

ACTION=$1
shift

case $ACTION in
  add)
    # shellcheck disable=SC2006
    uuid=`cat /proc/sys/kernel/random/uuid`
    echo "${uuid}=$1" >> "${CONFIG_FILE}"
  ;;
  del)
    line=0
    find=0
    while read -r LINE
        do
           line=$((line+1))
            if expr "$LINE" : '#' > /dev/null; then
              continue
            else
              # shellcheck disable=SC2092
              # shellcheck disable=SC2006
              NAME=`echo "${LINE}" | cut -d '=' -f 2`
              if [ "${NAME}" == "$1" ]; then
                  # shellcheck disable=SC2034
                  find=1
                  break
              fi
            fi
        done < "${CONFIG_FILE}"
    if [[ $line -gt 0 && ${find} -gt 0 ]]; then
        sed -i "${line}d" "${CONFIG_FILE}"
    fi
  ;;
    get)
        while read -r LINE
            do
                if expr "$LINE" : '#' > /dev/null; then
                    continue
                else
                    # shellcheck disable=SC2092
                    # shellcheck disable=SC2006
                    NAME=`echo "${LINE}" | cut -d '=' -f 2`
                    if [ "${NAME}" == "$1" ]; then
                            # shellcheck disable=SC2034
                            echo "${LINE}" | cut -d '=' -f 1
                            break
                    fi
                fi
            done < "${CONFIG_FILE}"
  ;;
esac
