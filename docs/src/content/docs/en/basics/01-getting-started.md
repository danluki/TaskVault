---
title: Getting started
description: Getting started.
---
## Introduction

Syncra is a distributed core(simple key-value storage for beginners) that is powered using your favourites technologies like [Raft](https://raft.github.io/) and Golang.

It provides a simple way to create a simple scalable stateful application based on your needs.

Syncra cluster have a leader, which is responsible for doing main job in cluster. If you familiar with distributed system you might think: why should I use single leader setup to instead of MultiRaft. And the answer is: simplicity. With Raft you can achive more than 100 nodes easily, and we will handle even datacenter fail for you. 

## State storage

Syncra is just a single binary, it stores the state in an internal BuntDB and replicate all changes between all server nodes using the Raft protocol, it doesn't need any other storage system outside itself.

## Configuration

See the [configuration](/en/basics/02-configuration/).

## Installation

See the [installation](/en/basics/03-installation/).

## Usage

Syncra by default uses the following ports:

- `8946` for serf layer between agents
- `8080` for HTTP for the API and Dashboard
- `6868` for gRPC and raft layer communication between agents.

