# Contributing to quic

## Setting up pre-commit hooks

### Via Docker

1. Install Docker if you don't have it already.

2. Run the pre-commit hooks using Docker:
   ```sh
   docker run --rm -v $(pwd):/code -w /code pre-commit run --all-files
   ```

3. To install the pre-commit hooks:
   ```sh
   docker run --rm -v $(pwd):/code -w /code pre-commit install
   ```

### Via Python

1. Install `pre-commit`:
   ```sh
   pip install pre-commit
   ```

2. Install the pre-commit hooks:
   ```sh
   pre-commit install
   ```

3. Run the hooks manually on all files:
   ```sh
   pre-commit run --all-files
   ```