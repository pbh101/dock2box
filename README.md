# In development

Currently we're focusing on migrating to a new backend. Once this is finished we should have a more complete version of this project.

# What?

Dock2Box is a tool that allows you to apply Docker or File images to bare-metal hardware using PXE boot.

# Why?

## Late vs early binding

Traditional host provisioning uses tools such as Cobbler, Kickstart and Configuration Management tools such as SaltStack, Ansible, Puppet or Chef to
finalize the host configuration. This means most of the complexity is late-binding i.e. while the host is being provisioned.

By moving the complexity to the build step on a CI server, you move most of the complexity where it will fail earlier rather than later.

Obviously not everything can be early-binding, but you can keep the things that are late-binding to a minimum and thus uncover errors sooner.

## Faster

At our current setup we can provision a fully configured server in less than 3 min. The main limitation is bandwidth and the size of the host image.

## Idempotent (repeatable, predictable)

Since we're dealing with a host image it's much more repeatable and predictable.

## Testable

One of the main issues with Configuration Management is testability: it's very hard to test since each deployment is slightly different.

With images they are basically immutable i.e. easier to test.

## DevOps

DevOps is all about breaking the barriers in-between Dev and Ops. This starts by adopting the same toolchain.

In using Docker for host provisioning you can use the same toolchain for software deployment and host provisioning.

## Running software on Bare-Metal

There are situations when you want to run software on bare-metal, but keep the deployment process close to what you
already do with Docker. Now it's just a matter of changing your Base Image "FROM ..." in your Dockerfile and rebuild.

# Workflow

![Workflow](img/workflow.png?raw=true)

# Tools

The backend and CLI is written in Go, the frontend will be written in AngularJS. Most of the provisioning scripts are for obvious reasons written in Bash.

There is a client library for Go and one in development for Python.

# Overview

![Overview](img/overview.png?raw=true)

# PXE menu

![PXE menu](img/pxe_menu.png?raw=true)
