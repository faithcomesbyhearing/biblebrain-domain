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

Initially, I am going to hack out a solution to removing S3 content that is not in biblebrain metadata. Once that works, I'll refactor to get it organized more logically

FIXME: cli shouldn't be in biblebrain-core... it's a client, but I don't want to mess with the go import of a dev module yet