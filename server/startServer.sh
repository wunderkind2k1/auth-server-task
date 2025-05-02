#! /usr/bin/env bash
JWT_SIGNATURE_KEY=$(cat ../keytool/keys/cddcbf9fe23b31ad.private.pem) go run main.go
