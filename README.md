# Consym



## Usage

### Choose a git branch

```bash
cd fabric

git checkout vanilla
# git checkout strawman # i.e. fabric-S
# git checkout consym

cd ..
```

### Build images

```bash
./consym.sh build
```

Build HLF and tape images

You must properly set these variables in `env.sh` in the root directory.

* BUILD_HLF_IMAGE
* BUILD_TAPE_IMAGE
* HLF_TYPE
* HLF_IMAGE_PREFIX

### Test

```bash
./consym.sh test
```

Setup network and test the performance of consym

You must properly set these variables in `env.sh` in the `test` directory.

* DEPLOY_MODE

### Clean

```bash
./consym.sh clean
```

Clean the images and/or test result.

You must properly set these variables in `env.sh` in the root directory.

* CLEAN_HLF_PEER_IMAGE
* CLEAN_TAPE_IMAGE
* CLEAN_RESULT
