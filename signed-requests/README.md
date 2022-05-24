# signedrequests

SignedRequest is written in Go and uses only Go in built packages. This tool for working with signedrequests urls and it provides users with following features

1. Generates Public and Private Keys using Ed25519 Algorirthm
2. Generate a base64 format for Public Keys
3. Generate a signed request url with signature


## Installation
Clone this repo. Change directory to gcp-data-drive.
```bash
git clone https://github.com/GoogleCloudPlatform/DIY-Tools
cd DIY-Tools/signed-requests/
```

## How to use?

To build `signedrequests` tool from script. For linux based environments, the following command will generate a runnable script name `signedrequests` in the same location.

```bash
bash-3.2$ go build signedrequests.go
```

Alternatively one can use `go run singedrequests.go` to work directly with the script itself.


#### Help

```bash
bash-3.2$ ./signedrequests --help
Usage: ./signedrequests subcommand [subcommand args...]

where: subcommand is one of generate-keys, encode-key, sign-url, sign-prefix, help
```

#### KeyPairs

To Ed25519 generate key pairs, use following command
```bash
./signedrequests generate-keys
```
̦
On succesful executions, `generate-keys` generated two files `private.key` and `public.pub` in the same location.

For more details, once can check `./signedrequests generate-keys -h`. As mentioned above, alternatively `go run signedrequests.go generate-keys -h` also works.

```bash
bash-3.2$ ./signedrequests generate-keys --help
Usage of generate-keys:
  -key string
    	file name into which to write the generated private key (default "private.key")
  -pub string
    	file name into which to write the generated public key (default "public.pub")
bash-3.2$

```

To generate a encoded public key from `public.pub`. By default, this command checks for the `public.pub` file in current locations. Check `./signedrequests encode-key --help` for more details
```bash
bash-3.2$ ./signedrequests encode-key
3j-vSlCYbpEViynAdfF4FXqG5csXqObthlEiUxKYBY8
```

#### Signed Request URL

To generate a signed request url,

```bash
bash-3.2$ ./signedrequests sign-url --url "http://example.com/bucket/sunny.jpg" --keyset example_keyset
http://example.com/bucket/sunny.jpg?Expires=1650542063&KeyName=example_keyset&Signature=3NOwNjha5yttobSKg3BveEzeBzTXQA3l4kij4QBJ0MT9I8QYrZQ0Sqyag4rwgU4_VcQTOaMmfQl0v9FOsmyqDw
```
The default time validity included is 1hr. For more details check `./signedrequests sign-url -h`



## Contributions


For contributions, please check [here](../contributing.md)
