The certificates used here were created using [minica](https://github.com/jsha/minica).  The following is a brief guide indicating how the certificates used here were created.

```
# starting from the project directory
cd CA
go get github.com/jsha/minica
minica --domains localhost --ip-addresses 127.0.0.1,::1
mv localhost/* ../static/certs
rmdir localhost
```

Trust the minica.pem public certificate for your CA in your browser so that it will trust your certificates.