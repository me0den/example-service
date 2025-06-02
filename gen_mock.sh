#!/bin/sh

mockery --output=app/api/v1/v1impl/mock --outpkg=mock --dir=app/api/v1/ --case=snake --name=.*Service$
mockery --output=app/api/v1/v1impl/mock --outpkg=mock --dir=domain/repo/ --case=snake --name=.*Repo$
