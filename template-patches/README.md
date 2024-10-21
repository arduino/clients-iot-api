## How to generate a patch for a generator

Patch has to be created starting from generator's git repo.
Read this guide for more info: https://openapi-generator.tech/docs/templating/

* clone openapi https://github.com/OpenAPITools/openapi-generator
* move to generator module and modify related mustache templata (for example, this the the path where go templates are stored: ``openapi-generator/modules/openapi-generator/src/main/resources/go/go.mod.mustache``)
* modify file
* move to modules path: ``openapi-generator/modules/openapi-generator/src/main/resources/``
* create diff file with command: ``git diff --relative > /path/to/diff-file.patch``
* 