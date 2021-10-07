# blockchain
My first own blockchain

## Resources

* See tutorial [Code your own blockchain in less than 200 lines of Go!](https://mycoralhealth.medium.com/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)
* See original code [nosequeldeebee/blockchain-tutorial](https://github.com/nosequeldeebee/blockchain-tutorial)
* Block hashing algorithm [Block_hashing_algorithm](https://en.bitcoin.it/wiki/Block_hashing_algorithm)
* Block hashing algorithm [o-create-a-bitcoin-hd-wallet-with-golang-and-grpc-part-l-u51d3wwm](https://hackernoon.com/how-to-create-a-bitcoin-hd-wallet-with-golang-and-grpc-part-l-u51d3wwm)
* Building a blockChain with persistence [building-a-blockchain-in-go-pt-iii-persistence](https://dev.to/nheindev/building-a-blockchain-in-go-pt-iii-persistence-3884)

## Command

- `make local` - launch the app
- `make local-proof` - launch the for pprof of work
- `make local-p2p` - launch the app for p2p
- `make local-proof-stake` - launch the app for pprof of stake
- `make local-seed` - launch the app for seed (key generator)
- `make local-persistence` - launch the app with local-persistance

### Metrics

Metrics can be get on http://127.0.0.1:8082/metrics depending on configuration

### Healthz

Status of server can be get by url

`
http://127.0.0.1:8082/liveness
`

and

`
http://127.0.0.1:8082/readiness
`

## Utiles

## Création et export d'une clé

```
// Création d'une clé publique
gpg --gen-key

// Export de la clé publique
gpg --export --armor adrienparrochia@gmail.com > pubkey.asc

scp -r -p pubkey.asc ariden@51.15.171.142:/home/ariden/

// Export de la clé publique
gpg --import pubkey.asc
```

## IPFS

## Ressources

* See [ipfs instanciate daemon](https://developers.cloudflare.com/distributed-web/ipfs-gateway/setting-up-a-server)
* See [ipfs tutorial](https://gist.github.com/YannBouyeron/53e6d67782dcff5995754b0a7333fa0b)
* [learn-to-securely-share-files-on-the-blockchain-with-ipfs](https://mycoralhealth.medium.com/learn-to-securely-share-files-on-the-blockchain-with-ipfs-219ee47df54c)
* [download ipfs](https://dist.ipfs.io/#ipfs-update)

### Installation d'IPFS

- Download [ipfs](https://dist.ipfs.io/#ipfs-update)

```
ipfs-update versions
ipfs-update install latest
ipfs init
sysctl -w net.core.rmem_max=2500000
ipfs daemon
```

### IPFS Daemon

To do this, we create a unit file at /etc/systemd/system/ipfs.service with the contents:

```
[Unit]
Description=IPFS Daemon

[Service]
ExecStart=/usr/local/bin/ipfs daemon
User=ipfs
Restart=always
LimitNOFILE=10240

[Install]
WantedBy=multi-user.target
```

### IPFS site perso

[Site de test](https://ipfs.io/ipfs/QmeY4kWRSpJUAseeeYet2AY4iCTT4G9DjQqhgEmRtA4q2D)


## TCP 

## Ressources

* See [networking tutorial](https://mycoralhealth.medium.com/part-2-networking-code-your-own-blockchain-in-less-than-200-lines-of-go-17fe1dad46e1)

## Test


In the first terminal :

```
make local-networking
```

In a second terminal : 

```
nc localhost 9000
5
7
...
```

