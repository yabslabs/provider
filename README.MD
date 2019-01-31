# Yabslabs github exporter

## Environment variables

### Basic auth

* GIT_ACCESS_TOKEN: personal access token
* GIT_USERNAME: personal username

## Steps

1. get all repos
2. create backup for all repos
3. waits until backup finnished
4. downloads backup
5. sends backup to defined storage module