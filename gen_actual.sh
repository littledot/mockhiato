#!/bin/bash
set -ev
mockhiato generate -p lib/generate/integration_test
mockhiato generate -p lib/generate/test/1.7/dependent_interface/actual
mockhiato generate -p lib/generate/test/1.7/nested_interface/actual
mockhiato generate -p lib/generate/test/1.7/same_import_name/actual
mockhiato generate -p lib/generate/test/1.7/sanity/actual
