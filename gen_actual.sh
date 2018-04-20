#!/bin/bash
set -ev
mockhiato generate -p lib/generate/integration_test
mockhiato generate -p lib/generate/test/dependent_interface/actual
mockhiato generate -p lib/generate/test/nested_interface/actual
mockhiato generate -p lib/generate/test/same_import_name/actual
mockhiato generate -p lib/generate/test/sanity/actual
