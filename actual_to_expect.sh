#!/bin/bash
set -ev
rm -rf lib/generate/test/1.7/dependent_interface/expect/
cp -rf lib/generate/test/1.7/dependent_interface/actual/ lib/generate/test/1.7/dependent_interface/expect/
rm -rf lib/generate/test/1.7/nested_interface/expect/
cp -rf lib/generate/test/1.7/nested_interface/actual/ lib/generate/test/1.7/nested_interface/expect/
rm -rf lib/generate/test/1.7/same_import_name/expect/
cp -rf lib/generate/test/1.7/same_import_name/actual/ lib/generate/test/1.7/same_import_name/expect/
rm -rf lib/generate/test/1.7/sanity/expect/
cp -rf lib/generate/test/1.7/sanity/actual/ lib/generate/test/1.7/sanity/expect/
