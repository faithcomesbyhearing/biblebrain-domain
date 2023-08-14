## biblebrain-domain

### Overview

The `biblebrain-domai` is a core Golang module for the BibleBrain project. This repository contains all essential domain models and repositories. It's designed as a central module, which means other BibleBrain repositories can easily import and utilize its Models and Repositories.

### Features

- Domain Models: Comprehensive structures that represent the different entities within the BibleBrain ecosystem.
- Repositories: Abstractions over data access methods, providing a unified way to query and interact with data sources.


### Getting Started

- Installing

Make sure you have Go installed and set up properly.
To include `biblebrain-domain` in you project:

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
