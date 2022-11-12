# progetto_SDCC

Pregoetto realizzato per il corso di Sistemi Distribuiti e Cloud Computing, Ingegneria Informatica LM Università Roma Tor Vergata.

Il progetto simula l'esecuzione di tre algoritmi di mutua esclusione distribuita ( token centralizzato, token decentralizzato, ricart & agrawala ).

## Utilizzo

### Prima di iniziare

Per eseguire questa applicazione è prevista l'installazione di [Docker](https://www.docker.com/) e di [Docker compose](https://docs.docker.com/desktop/install/linux-install/)

La directory [mutual_exclusion](https://github.com/MegAtonBoom/progetto_SDCC/tree/main/mutual_exclusion) prevede tutto il necessario per il corretto avvio dell'applicazione.

### Linux
Utilizzare lo script [mes.sh](https://github.com/MegAtonBoom/progetto_SDCC/blob/main/mutual_exclusion/mes.sh) per l'avvio. Lo script prevede diversi flag di configurazione:
- "-n"    specifica numero di peer
- "-v"    *verbose*: le attività del processo verranno loggate in una apposita directory "logs" 
- "-d"    *delay*: tutte le chiamate a rpc subiranno un delay variabile, a simulare una congestione di rete
- "-a"    specifica algoritmo da utilizzare per il processo, nello specifico:
  - 'c'     token centralizzato
  - 'd'     token decentralizzato 
  - 'r'     ricart & agrawala
- "-h"    display messaggio di aiuto con indicazioni su come eseguire correttamente lo script

**N.B.:** L'avvio dell'aplicazione è previsto solo tramite l'esecuzione di questo script.

### Windows ed altri
Non previsto.
