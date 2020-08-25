This package allows us to depend on other modules, to _force_ these other
modules to upgrade to some minimum version. This allows us to express "co
dependency" requirements in cases where this module doesn't strictly speaking
_depend_ on another module, but conflicts with some version of that module.
We are using this here to depend on deprecated modules that have been
merged into this package. 

In practice, this means:

1. Packages imported here _will not_ end up in the final binary as nothing
   imports this package.
2. Modules containing these packages will, unfortunately, be downloaded as
   "dependencies" of this package.
3. Anyone using this module will be forced to upgrade all co-dependencies.
