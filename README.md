# cdkctl
 
This tool is benefical if you intend to run cdk stacks concurrently(parallel) thus saving time to run each stack sequentially.

## Contents
* [Installing](#Installing)
    * Install using `make`
    * Install using `brew`
* [Commands](#Commands)

## Installing

Use the following to be able to install on MacOS or Linux:

### Install using `make`

1. Clone the repo
2. Make sure you have go > go1.14.2 installed
3. Run `make local`

### Install using `brew`

```bash
 brew tap maheshrayas/cdkctl
 ```

 ```bash
  brew install cdkctl
  ```

### Install on Windows & Linux

* Download the latest release from https://github.com/maheshrayas/cdkctl/releases/tag/v.0.2.0
* Unpack and set the binary in the `PATH`
* For windows, rename the binary to `cdkctl.exe`

## Commands

### Run all the stacks in parallel

Refer: [stacks.json](./example-config/stacks-all.json) and [args.json](./example-config/args.json) to describe stacks name and context arguments (runtime)

Sample commands:
1. With no run time arguments

```bash
cdkctl deploy --stacks-file configs/stacks.json --tool-kit toolkit-name
```

2.With runtime (context) arguments

```bash
cdkctl deploy --stacks-file example-configs/stacks-all.json --tool-kit toolkit-name --args example-configs/args.json
```

3.If stacks are dependent on each other you can frame the json as described in [stacks-dependent](./example-config/stacks-dependent.json)

```bash
cdkctl deploy --stacks-file example-configs/stacks-dependent.json --tool-kit toolkit-name --args example-configs/args.json
```

4.Destroy all the stacks

```bash
cdkctl destroy --stacks-file example-configs/stacks-all.json --tool-kit toolkit-name
```
