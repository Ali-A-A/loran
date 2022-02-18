# loran

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ali-a-a/loran/ci?label=ci&logo=github&style=flat-square)
[![Codecov](https://img.shields.io/codecov/c/gh/ali-a-a/loran?logo=codecov&style=flat-square)](https://codecov.io/gh/ali-a-a/loran)

A distinct counter using the HyperLogLog algorithm.

## Description

Sometimes, we need to count the distinct count of an entity that
is observed by some users. \
Using the HyperLogLog algorithm we could approximate the number of distinct elements in a multiset.
For better performance, this service has three components, i.e., Abacus, Cranmer, and Scheduler. For 
more information, you can see the architecture section.
