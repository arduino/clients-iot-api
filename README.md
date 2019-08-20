# iot-api clients generator

This repo contains the informations and the tools needed to automatically
generate API clients for the `iot-api` service in [any language][0] supported by
OpenAPI generator.

The OpenAPI generator is orchestrated by a Python tool called [apigentools][1]
that let us keep the spec files and the configuration bits in one single repo
(this one) without duplicating the boilerplate on each git repo hosting the
actual clients.

## Requirements

* Python 3.6+
* OpenAPI generator 4.0+
* Apigentools (`pip install apigentools`)

## Workflow

Apart from the initial setup (the procedures needed to setup this git repo),
the ideal workflow consists in updating this repo every time the API service
changes, re-generate all the clients, push the generated code to the git repo
it belongs to and release the updated clients (this last step may vary
depending on the programming language).

The operations are detailed in the following paragraphs.

### Get an updated version of the API specification

In this case the specs are generated by Goa using Swagger and they can be found
at [http://api-dev.arduino.cc/iot/swagger.json][2]. To be more future proof
and to leverage the latest versions of the tools available, we're using the
version 3 of the OpenAPI for the generator, this means that the output from
Goa must be converted into a compatible spec [like the one in this repo][3].

Several conversion tools from OpenAPI v2 to v3 exist, and some of them can be
easily integrated in a CI, so major updates will be performed by either:

* have Goa produce a v3 OpenAPI spec
* dump the Goa spec, convert to v3 and update this repo

Minor updates might be done manually since v3 uses Yaml and the resulting spec
is human friendly.

### Validate the Spec

The generator can validate the content of the spec, this should likely go in a
CI step:

```sh
apigentools validate
```

### Generate the clients

If validation succeeded, chances are that this step is as easy as running:

```sh
apigentools generate --builtin-templates
```

`apigentools` has a simple script that can be invoked for each client generated
that will push the resulting code to different repos. A release process should
take from there and shouldn't be part of this workflow.

[0]: https://openapi-generator.tech/docs/generators
[1]: https://github.com/DataDog/apigentools
[2]: http://api-dev.arduino.cc/iot/swagger.json
[3]: clients-iot-api/spec/v2/swagger.yaml
