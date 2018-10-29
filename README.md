# Epitome (formerly Pond)

_Current version:_ `0.1.1`

## Table of Contents

1. [Introduction & Overview](#introduction--overview)
2. [Epitome as a Standard](#epitome-as-a-protocol)
3. [Epitome as an App](#epitome-as-an-app)
   * [Installing](#installing)
   * [Building](#building)
   * [Usage](#usage)
   * [What's Missing for 1.0?](#whats-missing-for-10)
4. [Information for Contributors](#information-for-contributors)
   * [Forking](#forking)
   * [Testing, Bug Reporting & Feature Requests](#testing-bug-reporting--feature-requests)
   * [Protocol Definition](#protocol-definition)


## Introduction & Overview

Epitome is a protocol that aims to standardize Atom+RSS syncing across the web. It's also a self-hosted alternative to the backend part of Google Reader, written in Go.

## Epitome as a Standard

Epitome defines a set of [RESTful HTTP API endpoints](https://github.com/ArturoVM/epitome/blob/master/api_doc.md) that receive and send data in a standardized format, while making it very comfortable for developers to use. Comfort, ease of use, minimalism and elegance are Epitome's primary design goals.

Anybody can implement Epitome using their technology stack of choice. It could be used in a multi-user environment (as is the case with the reference implementation), or a single-user environment.

Epitome also aims to allow developers to build upon it: other protocols can be thrown on top, as long as the implementation satisfies all the API endpoints.

All of this allows for a (hopefully) bright future: if developers were to adopt Epitome in their services, maximum portability would be achieved, without all the pain that comes with learning and implementing new APIs; let alone the mess that RSS is.

## Epitome as an App

The Epitome reference implementation does _not_ include a feed reader; it is just a backend that feed reading apps can use to sync user's feeds and articles across devices. 

It's worth noting that I aim to make the Epitome reference implementation as high quality as possible so it can also be used in production, and it's not just there as a learning tool.

### Installing

Because Go compiles and links all libraries statically, you can download the appropriate binary for your platform, and run it right away. However, you can of course always build it from source.

### Building

To build Epitome, you need the following dependencies:

* moovweb/gokogiri
* bmizerany/pat
* mattn/go-sqlite3
* robfig/cron
* code.google.com/p/go.crypto/bcrypt

```bash
go get -u github.com/moovweb/gokogiri
go get -u github.com/bmizerany/pat
go get -u github.com/mattn/go-sqlite3
go get -u github.com/robfig/cron
go get -u code.google.com/p/go.crypto/bcrypt
```

Then, do:

```bash
git clone git://github.com/ArturoVM/epitome.git
cd epitome
go build -o epitome *.go
```

If you want cutting edge, before building, do:

```bash
git checkout develop
```

### Usage

To run the Epitome binary, simply do:

```bash
./epitome
```

To specify a port other than the default one, use the `--port` flag (or its shorthand `-p`):

```bash
./epitome -p 8080
```

To enable verbose mode (useful for debugging and development), use the `--verbose` flag:

```bash
./epitome -p 8080 --verbose
```

During verbose mode, you can choose whether to log requests' bodies or not (default is off):

```bash
./epitome -p 8080 --verbose --log-body
```

If verbose mode is off, this option is a no-op. Keep in mind that if you use this option during verbose mode, memory usage could increase significantly (definitely not recommended during production).

### What's Missing for 1.0?

#### First and foremost? Testing, testing, testing.

Atom is amazing and awesome and unicorns and rainbows. RSS is not. RSS is a messy pain in the ass—kinda like diarrhea. I need people to test the program and try to break it (bug reporting will be _very_ appreciated), so I can sort out all the little quirks and handle all sorts of use cases, and make the Epitome reference implementation much more robust over time.

#### OPML Import and Export

Which will very likely be added in 0.2.0

#### Favorites API

Which will also very likely be added in a near-future minor version.

#### Support for Enclosures

This involves a JSON schema addition. Might be tackled on 0.2.x

#### That's it?

No. I'm sure more things will be added to this section as time goes by.

## Information for Contributors

You can contribute in any of the following ways:

### Forking

To contribute with a bug fix or to add a feature, fork the repo, create a new branch and add a pull request. If you're fixing a bug, name the branch after the issue ID.

You can very easily infer what's in each file based on its name (e.g. route definitions are in `routes.go`, the articles API is in `articles.go`, etc).

### Testing, Bug Reporting & Feature Requests

Just download the app, use it as much as you can (but, of course, don't rely on it just yet) and submit any bugs you encounter or feature requests you may have.

### Standard Definition

You can contribute a lot just by giving your views of the standard. If there's something you think isn't quite right, or know how it could be improved, submit it to the issue tracker.
