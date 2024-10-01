## biblebrain-domain

### Overview

The `biblebrain-domain` is a core Golang module for the BibleBrain project. This repository contains all essential domain models and repositories. It's designed as a central module, which means other BibleBrain repositories can easily import and utilize its Models and Repositories.

### Features

- Domain Models: Comprehensive structures that represent the different entities within the BibleBrain ecosystem.
- Repositories: Abstractions over data access methods, providing a unified way to query and interact with data sources.


### Getting Started

- Installing

Make sure you have Go installed and set up properly.
To include `biblebrain-domain` in your project:

```bash
GOPRIVATE=github.com/faithcomesbyhearing
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

```bash
go get github.com/faithcomesbyhearing/biblebrain-common
```

- Usage

Once installed, you can import the module in your project:

```bash
import "github.com/faithcomesbyhearing/biblebrain-common"
```

Use the models and repositories in your services as needed.


### BWF notes
Preferred name is biblebrain-core instead of biblebrain-domain.

Maybe packages should be organized by type of domain (eg stocknumber, language, product, media)

Initially, I am going to build a solution to removing S3 content that is not in biblebrain metadata. Once that works, I'll refactor to get it organized more logically

a) create tunnel so this can be run locally
 ssh -v -N -L 3306:rds.dev.biblebrain.com:3306 ec2-user@dbp-dev-bastion
b) build
cd cli/lint
go build
go install

c) run audit
lint audit --bucket dbp-prod --filesetId BNGCLVN2DA 

if there are orphans, a file will be created in directory audited
after reviewing, move the file to toRemove 

d) run remove
lint remove --bucket dbp-prod --filesetId BNGCLVN2DA 

files will be removed, and the json file moved to directory "removed"

FIXME: cli shouldn't be in biblebrain-core... it's a client, but I don't want to mess with the go import of a dev module yet

SQLC usage: after editing queries.sql, run 
cd adapters/mysql && sqlc generate