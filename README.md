# blockchain

[![ariden83](https://codecov.io/gh/ariden83/blockchain/branch/master/graph/badge.svg)](https://codecov.io/gh/ariden83/blockchain)
[![Build Status](https://travis-ci.org/ariden83/blockchain.svg?branch=master)](https://travis-ci.org/ariden83/blockchain)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fariden83%2Fblockchain.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fariden83%2Fblockchain?ref=badge_shield)

New blockchain. See website of project [blockchain-altcoin.com](https://www.blockchain-altcoin.com/)

## Resources

### Blockchain 

* See tutorial [tutorial](https://mycoralhealth.medium.com/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)
* See original code [tutorial](https://github.com/nosequeldeebee/blockchain-tutorial)
* Block hashing algorithm [tutorial](https://en.bitcoin.it/wiki/Block_hashing_algorithm)
* Create a bitcoin hd wallet [tutorial](https://hackernoon.com/how-to-create-a-bitcoin-hd-wallet-with-golang-and-grpc-part-l-u51d3wwm)
* Building a blockChain with persistence [tutorial](https://dev.to/nheindev/building-a-blockchain-in-go-pt-iii-persistence-3884)
| [code source](https://replit.com/@nheingit/GolangBlockChain-3)
* Building a blockChain with transactions [tutorial](https://dev.to/nheindev/building-a-blockchain-in-go-pt-iv-transactions-1612)
| [code source](https://replit.com/@nheingit/GolangBlockChain-4)
* Ethereum block structure explained [tutorial](https://medium.com/@eiki1212/ethereum-block-structure-explained-1893bb226bd6)
* Ethereum accounts transactions gas [tutorial](https://hudsonjameson.com/2017-06-27-accounts-transactions-gas-ethereum/)
* Le minage en 7 étapes [tutorial](http://www.ab-consulting.fr/blog/blockchain/minage-7-etapes=)
* libp2p in go tutorial [tutorial](https://dev.to/feliperosa/getting-started-with-libp2p-in-go-4hoa)
* p2p discovery mdns [code source](https://gitlab.dms3.io/p2p/go-p2p/-/blob/master/p2p/discovery/mdns.go)
* p2P examples [examples](https://github.com/libp2p/go-libp2p/tree/master/examples/)
* blockchain go project by Nomad [code source](https://github.com/librity/nc_nomadcoin)
* badger - database [tutorial](https://dgraph.io/docs/badger/get-started/)
* cipher GCM [tutorial](https://pilabor.com/blog/2021/05/js-gcm-encrypt-dotnet-decrypt/)
* Clés privées, clés publiques et adresses dans Bitcoin [tutorial](https://cryptoast.fr/cles-privees-cles-publiques-et-adresses-dans-bitcoin/)
* Comment marchent les transactions bitcoin [tutorial](https://www.pensezblockchain.ca/les-transactions-bitcoin-partie-1)
* Qu'est-ce que P2PKH [tutorial](https://academy.bit2me.com/fr/qu%27est-ce-que-p2pkh/)
* Comment fonctionne P2PKH [tutorial](https://learnmeabitcoin-com.translate.goog/technical/p2pkh?_x_tr_sl=auto&_x_tr_tl=fr&_x_tr_hl=fr&_x_tr_pto=wapp)
* Créez une transaction Bitcoin brute et signez-la avec Golang [code source](https://ichi.pro/fr/creez-une-transaction-bitcoin-brute-et-signez-la-avec-golang-165707908919466)
* btcd is an alternative full node bitcoin implementation written in Go (golang) [code source](https://github.com/btcsuite/btcd)
* Implementing RSA Encryption and Signing in Golang [tutorial](https://www.sohamkamani.com/golang/rsa-encryption/) [code source](https://gist.github.com/sohamkamani/08377222d5e3e6bc130827f83b0c073e)
* Decred is a blockchain-based cryptocurrency [code source](https://github.com/decred/dcrd)
* Recommended practices for secure signature generation [tutorial](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/2019-05-15-schnorr.md#recommended-practices-for-secure-signature-generation)
* Quantum computation simulator  [code source](https://github.com/itsubaki/q)

## Keys Encrypting

1. Ce qui est permanent
- user : priv key
- blockchain : pub key et pub blockchain key

2. actions
- user generate priv blockchain key with priv key
- priv key generate pub blockchain key and a script
- blockchain validator verify script is valid with pub blockhain key (possibilité de vérifier que script decode avec blockchain key = public key)

3. PB: 
- (pas de vérification si pub blockchain key = public key) > peut être résolu 
- (pas si c'est bien la priv key de l'utilisateur qui fait la transaction pour les validators)

## Installation

- Code verification with [code source](https://github.com/securego/gosec)

## Command

- `make local` - launch the app
- `local-light` - launch the app without files generators (logs, blockchain, ...)

### Build

- `make proto` - generate proto files
- `make local-vendor` - generate vendor files

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

### Features

- [x] Récupère l'ensemble des seeds et la blockchain complète lors de la conexion avec le premier server
- [x] Generate a new seed
- [x] Access to your wallet
- [x] Mine a new bloc
- [x] Send tokens to another
- [x] Access to your balance
- [x] Les blocs créés sont obligatoirement validés par plus de 50% des serveurs connectés
- [x] Le service requiert minimum deux serveurs pour fonctionner
- [x] Liste les serveurs actuellement actifs
- [x] Ajout des frais de transactions en faveur du mineur
- [x] Evolution de la difficulté

### Future
- [ ] Fully Tested
- [ ] Add bench
- [ ] GRPC endpoint
- [ ] Update seed database
- [ ] Create seed database with seed database on network
- [ ] Add oauth
- [ ] Encrypt data with cypher key
- [ ] Load seed database in many times
- [ ] Load blockchain database in many times
- [ ] Ajouter un champ metadata dans chaque seed (pour infos complementaires)
- [ ] Se brancher à l'API metamask
- ...

### Godocs
You can browse the documentation of all local packages and projects with the Godocs package:

```
go install golang.org/x/tools/godoc
godoc -http=:6060
```

This will install the executable and start a server listening on http://localhost:6060

### Test

#### 1) Open services

You have to open two terminal minimun.

In the first terminal :

```
make local
```

return 

```
go run main.go -p2p_target /ip4/127.0.0.1/tcp/8098/p2p/QmWV1qKRBSy8vggYgMSWDGukmwcus8wbuSoru31oNaEWdd
```

Then launch one or more light service * 

```
cd example/light
go run main.go -p2p_target /ip4/127.0.0.1/tcp/8098/p2p/QmWV1qKRBSy8vggYgMSWDGukmwcus8wbuSoru31oNaEWdd
```

* light service dont create files, it's just for tests

#### 2) GENERATE A SEED (WALLET)

```
make local
```

And call IT

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://127.0.0.1:8098/wallet' -d ''
```

Return 

```
{
  "Address": "1P1aBegXRiTinJhhEYHHiMALfG26Wu9sG3",
  "Timestamp": "2021-10-11 16:52:12.416519751 +0200 CEST m=+27.320229089",
  "PubKey": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF",
  "PrivKey": "xprv9s21ZrQH143K2yV5cD5SYvTzUP1XyXgGVWd1FLWY3kxTbAUQLHqtcMnJXJgfkH1Q3UqXqZ6FmDRTwLHdvDTJC6wNm7Vh9FokRma8WrDGQAe",
  "Mnemonic": "couple robot escape silent main once smoke check good basket mimic similar"
}
```


#### 3) GENERATE YOUR FIRST BLOCK 

And call IT in the first terminal

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://127.0.0.1:8098/write' -d '{"address": "1P1aBegXRiTinJhhEYHHiMALfG26Wu9sG3", "key": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF"}'
```

Return 

```
{
  "Index": 26,
  "Timestamp": "2021-10-11 17:04:46.977307004 +0200 CEST m=+45.479261977",
  "Transactions": [
    {
      "ID": null,
      "Inputs": [
        {
          "ID": "",
          "Out": -1,
          "Sig": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF"
        }
      ],
      "Outputs": [
        {
          "Value": 1,
          "PubKey": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF"
        }
      ]
    }
  ],
  "Hash": "MGRiODRmMWFlNjhmZjQ5ZDA5ZmI4M2JhODE0MDg2YTdjN2QxOWYyZGFjODEzMzdhZmVlMTU3YjU4MjZhYzkwZQ==",
  "PrevHash": "MDcyYWMxYTlkNmI5YjQ1ZWFiMWYyMTI3Y2U1YzVlMGVlZjBiYjE3NTI3NTFkNzQyMWM2Y2U1ZmUxN2MwOTUyNA==",
  "Difficulty": 1,
  "Nonce": "8"
}
```

#### 4) YOUR BALANCE

Call your balance :

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://127.0.0.1:8098/balance' -d '{"key": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF"}'
```

return 

```
Balance of xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF: 1
```

#### 5) GENERATE A SECOND WALLET 

```
{
  "Address": "1NKEsiake5Yu8yx2H2uHm2oJZe2xYnQ8ZS",
  "Timestamp": "2021-10-11 17:04:24.263414034 +0200 CEST m=+3.378949045",
  "PubKey": "xpub661MyMwAqRbcG4VYfVo7ptRncn7wsGMjNubLNrm5Stu5ERP4RtJqo7sQgSQAESwyJKi442EJ6sNWRz5wWZ2ecFE8p1JEJs6qGkzPKncdkhb",
  "PrivKey": "xprv9s21ZrQH143K3aR5ZUG7TkV44kHTTodt1gfjaUMTtZN6Md3utLzbFKYvqCuqyXAnVcirzpNuzcBkcvpTfJNRjakAwsmEA26wNWmDmLJKXYD",
  "Mnemonic": "couple office mix shadow glide crater sister check gown sister mirror indoor"
}
```

#### 6) SEND ONE TOKEN TO THE 2nd WALLET

a) On envoi la somme du compte A au compte B

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://127.0.0.1:8098/send' -d '{"from": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF", "to": "xpub661MyMwAqRbcG4VYfVo7ptRncn7wsGMjNubLNrm5Stu5ERP4RtJqo7sQgSQAESwyJKi442EJ6sNWRz5wWZ2ecFE8p1JEJs6qGkzPKncdkhb", "amount": 3}'
```

b) On récupère la balance du compte envoyeur

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://127.0.0.1:8098/balance' -d '{"key": "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF"}'
```

retourne 


```
Balance of xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF: 97
```

c) On récupère la balance du compte receptionneur

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://127.0.0.1:8098/balance' -d '{"key": "xpub661MyMwAqRbcG4VYfVo7ptRncn7wsGMjNubLNrm5Stu5ERP4RtJqo7sQgSQAESwyJKi442EJ6sNWRz5wWZ2ecFE8p1JEJs6qGkzPKncdkhb"}'
```

retourne 


```
Balance of xpub661MyMwAqRbcG4VYfVo7ptRncn7wsGMjNubLNrm5Stu5ERP4RtJqo7sQgSQAESwyJKi442EJ6sNWRz5wWZ2ecFE8p1JEJs6qGkzPKncdkhb: 3
```

#### 7) Communicate new update of blockChain / wallet with every blockChain service

Après la commande 

```
make local
> 2021-10-15T16:40:38.669+0200	INFO	Now run "go run main.go -l 8198 -d /ip4/127.0.0.1/tcp/8097/p2p/QmdJboshgG8BuRexqmq9opEsr49Zw961UqSMQrrfXxyzxQ" on a different terminal
```

il faut récupérer l'adresse TCP transmise dans les logs et l'executer dans un nouveau terminal

```
cd ./cmd/p2p
go run main.go -l 8198 -d /ip4/127.0.0.1/tcp/8097/p2p/QmdJboshgG8BuRexqmq9opEsr49Zw961UqSMQrrfXxyzxQ
```

Après chaque création / update de la blockChain ou des seeds, le deuxième service lancé va se mettre à jour


![minage](https://github.com/ariden83/blockchain/blob/main/readme/minage.png)

![messaging](https://github.com/ariden83/blockchain/blob/main/readme/messages.png)

## GPG tutorial

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

## IPFS tutorial

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

> [Site de test](https://ipfs.io/ipfs/QmeY4kWRSpJUAseeeYet2AY4iCTT4G9DjQqhgEmRtA4q2D)


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


