# blockchain
My first own blockchain

## Resources

* See tutorial [Code your own blockchain in less than 200 lines of Go!](https://mycoralhealth.medium.com/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)
* See original code [nosequeldeebee/blockchain-tutorial](https://github.com/nosequeldeebee/blockchain-tutorial)
* Block hashing algorithm [Block_hashing_algorithm](https://en.bitcoin.it/wiki/Block_hashing_algorithm)


## Command

- `make local` - launch the app
- `make test` - launch test


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

#### Création d'une clé publique

```
gpg --gen-key
```

#### Export de la clé publique

```
gpg --export --armor -email > pubkey.asc
```

## Ressources

* Block hashing algorithm [learn-to-securely-share-files-on-the-blockchain-with-ipfs](https://mycoralhealth.medium.com/learn-to-securely-share-files-on-the-blockchain-with-ipfs-219ee47df54c)

