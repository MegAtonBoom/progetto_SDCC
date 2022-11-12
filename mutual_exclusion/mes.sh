#!/bin/bash

help() {
    echo "Usage: sh  ${0} [CONFIG FLAGS]"
    echo "Config:"
    echo  "FLAG   VALUES                    DESCRIPTION"
    echo  "-n     integer number>0          Number of peers to spawn"
    echo  "-a     'c' or 'd' or 'r'         Algorithm to use- c is for centralized token, d for decentralized, r for ricart & agrawala"
    echo  "-v                               Verbose mode- writes on log"
    echo  "-d                               Delay- with this flasg every rpc call will require some time, simulating internet congestion" 
    echo  "Defaults, if some of the configurations above are not provided:"
    echo  "Number of peers = 4, Algorithm = centralized token, Verbosity = false, Delay = false"
    echo  "A config flag with a non planned value will result in the defaul setting being applied"
}

NPEERS=4
VERBOSE="false"
DELAY="false"

COMPOSE_FILE="docker-compose1.yml"
ENV_FILE="./.env"
LOG_DIR="./logs"

# Parsing command line opts
while getopts "hvdn:a:" opt; do

    case ${opt} in
        h ) 
            help
            exit 0
            ;;
        n )
            NPEERS=$OPTARG   
            ;;
            
        v )
            VERBOSE="true"
            ;;
        a )
            if [ "${OPTARG}" == "d" ]; then
                COMPOSE_FILE="docker-compose2.yml"
            fi
            if [ "${OPTARG}" == "r" ]; then
                COMPOSE_FILE="docker-compose3.yml"
            fi
            ;;
        d )
            DELAY="true"
            ;;
        ? )
            help
            exit 1 
    esac
done

shift $((OPTIND -1))


echo  "Current Settings"
echo  "Peer number: ${NPEERS}"
echo  "Verbosity: ${VERBOSE}"
echo  "Delay: ${DELAY}"

# env file
rm -f ${ENV_FILE}
touch ${ENV_FILE}
echo "NPEERS = \"${NPEERS}\"" >> ${ENV_FILE}
echo "VERBOSE = \"${VERBOSE}\"" >> ${ENV_FILE}
echo "DELAY=\"${DELAY}\"" >> ${ENV_FILE}
echo "[+] Created environment file for docker-compose"

# Startup

sudo rm -d -r ${LOG_DIR}
mkdir ${LOG_DIR}
docker compose -f docker-compose1.yml -f docker-compose2.yml -f docker-compose3.yml stop
docker compose -f docker-compose1.yml -f docker-compose2.yml -f docker-compose3.yml down --remove-orphans
docker image rm -f mutual_exclusion-registration_dc
docker image rm -f mutual_exclusion-peer_dc
docker image rm -f mutual_exclusion-coord_dc
docker compose build
echo "docker compose -f ${COMPOSE_FILE} up --scale peer_dc=${NPEERS}"
docker compose -f ${COMPOSE_FILE}  up --scale peer_dc=${NPEERS}
