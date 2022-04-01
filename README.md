```
docker-compose up --build
```

```
# generate keys
cd keys
openssl genrsa -out id_rsa 4096
openssl rsa -in cert/id_rsa -pubout -out id_rsa.pub
# THERE MUST BE AN EXTRA NEWLINE AT THE END OF THE PUBLIC KEY
# base 64 to put into env vars
base64 id_rsa > private_key_b64
base64 id_rsa.pub > public_key_b64
```