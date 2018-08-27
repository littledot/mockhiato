#!/bin/bash
set -ev
mockhiato generate -p lib/generate/integration_test -m {func}Run
mockhiato generate -p lib/generate/testdata/1.7/dependent_interface/actual -m {func}Run
mockhiato generate -p lib/generate/testdata/1.7/nested_interface/actual -m {func}Run
mockhiato generate -p lib/generate/testdata/1.7/no_defs_only_uses/actual -m {func}Run
mockhiato generate -p lib/generate/testdata/1.7/same_import_name/actual -m {func}Run
mockhiato generate -p lib/generate/testdata/1.7/sanity/actual -m {func}Run
mockhiato generate -p lib/generate/testdata/1.9/type_alias/actual -m {func}Run

rm -rf lib/generate/testdata/1.7/dependent_interface/expect/
cp -rf lib/generate/testdata/1.7/dependent_interface/actual/ lib/generate/testdata/1.7/dependent_interface/expect/
rm -rf lib/generate/testdata/1.7/nested_interface/expect/
cp -rf lib/generate/testdata/1.7/nested_interface/actual/ lib/generate/testdata/1.7/nested_interface/expect/
rm -rf lib/generate/testdata/1.7/no_defs_only_uses/expect/
cp -rf lib/generate/testdata/1.7/no_defs_only_uses/actual/ lib/generate/testdata/1.7/no_defs_only_uses/expect/
rm -rf lib/generate/testdata/1.7/same_import_name/expect/
cp -rf lib/generate/testdata/1.7/same_import_name/actual/ lib/generate/testdata/1.7/same_import_name/expect/
rm -rf lib/generate/testdata/1.7/sanity/expect/
cp -rf lib/generate/testdata/1.7/sanity/actual/ lib/generate/testdata/1.7/sanity/expect/
rm -rf lib/generate/testdata/1.9/type_alias/expect/
cp -rf lib/generate/testdata/1.9/type_alias/actual/ lib/generate/testdata/1.9/type_alias/expect/
