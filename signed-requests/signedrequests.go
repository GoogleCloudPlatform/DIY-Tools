package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// genKeys generates a key pair suitable for request signing
func genKeys(priv, pub io.Writer) error {
	pubKey, privKey, err := ed25519.GenerateKey( /*rand=*/ nil)
	if err != nil {
		return fmt.Errorf("could not generate key: %v", err)
	}

	if _, err := priv.Write(privKey); err != nil {
		return fmt.Errorf("could not write private key: %v", err)
	}

	if _, err := pub.Write(pubKey); err != nil {
		return fmt.Errorf("could not write public key: %v", err)
	}

	return nil
}

// encodeKey base64 encodes the key for use with our configuration
func encodeKey(key []byte) {
	fmt.Fprintln(os.Stdout, base64.RawURLEncoding.EncodeToString(key))
}

// signURL signs the given path using query parameters
func signURL(key []byte, keyset, url string, expires time.Time) {
	sep := '?'
	if strings.ContainsRune(url, '?') {
		sep = '&'
	}
	toSign := fmt.Sprintf("%s%cExpires=%d&KeyName=%s", url, sep, expires.Unix(), keyset)
	sig := ed25519.Sign(key, []byte(toSign))
	fmt.Fprintf(os.Stdout, "%s&Signature=%s\n", toSign, base64.RawURLEncoding.EncodeToString(sig))
}

// signPrefixWithQueryParameters signs the given prefix using query parameters
func signPrefixWithQueryParameters(key []byte, keyset, prefix string, expires time.Time) {
	toSign := fmt.Sprintf("URLPrefix=%s&Expires=%d&KeyName=%s", base64.RawURLEncoding.EncodeToString([]byte(prefix)), expires.Unix(), keyset)
	sig := ed25519.Sign(key, []byte(toSign))
	fmt.Fprintf(os.Stdout, "%s&Signature=%s\n", toSign, base64.RawURLEncoding.EncodeToString(sig))
}

// signPrefixWithCookie signs the given prefix using a cookie
func signPrefixWithCookie(key []byte, keyset, prefix string, expires time.Time) {
	toSign := fmt.Sprintf("URLPrefix=%s:Expires=%d:KeyName=%s", base64.RawURLEncoding.EncodeToString([]byte(prefix)), expires.Unix(), keyset)
	sig := ed25519.Sign(key, []byte(toSign))
	fmt.Fprintf(os.Stdout, "Edge-Cache-Cookie=%s:Signature=%s\n", toSign, base64.RawURLEncoding.EncodeToString(sig))
}

// signPrefixWithPathComponent signs the given prefix using path components
func signPrefixWithPathComponent(key []byte, keyset, prefix string, expires time.Time) {
	// Remove trailing slashes because the path component format starts
	// with a slash and we don't want duplicated slashes in the URL.
	prefix = strings.TrimRight(prefix, "/")
	toSign := fmt.Sprintf("%s/edge-cache-token=Expires=%d&KeyName=%s", prefix, expires.Unix(), keyset)
	sig := ed25519.Sign(key, []byte(toSign))
	fmt.Fprintf(os.Stdout, "%s&Signature=%s/\n", toSign, base64.RawURLEncoding.EncodeToString(sig))
}

func main() {
	genCmd := flag.NewFlagSet("generate-keys", flag.ExitOnError)
	genKey := genCmd.String("key", "private.key", "file name into which to write the generated private key")
	genPub := genCmd.String("pub", "public.pub", "file name into which to write the generated public key")

	ekCmd := flag.NewFlagSet("encode-key", flag.ExitOnError)
	ekPub := ekCmd.String("pub", "public.pub", "file name from which to read the public key")

	suCmd := flag.NewFlagSet("sign-url", flag.ExitOnError)
	suKey := suCmd.String("key", "private.key", "file name from which to read the private key")
	suKeyset := suCmd.String("keyset", "", "the name of the EdgeCacheKeyset to use.  Must not be the empty string.")
	suURL := suCmd.String("url", "", "the URL to sign, including protocol.  Must not be the empty string.  For example: http://example.com/path/to/content")
	suTTL := suCmd.Duration("ttl", time.Hour, "duration the signed request is valid")

	spCmd := flag.NewFlagSet("sign-prefix", flag.ExitOnError)
	spKey := spCmd.String("key", "private.key", "file name from which to read the private key")
	spKeyset := spCmd.String("keyset", "", "the name of the EdgeCacheKeyset to use.  Must not be the empty string.")
	spURL := spCmd.String("url-prefix", "", "the URL prefix to sign, including protocol.  Must not be the exmpty string.  For example: http://example.com/path/ for URLs under /path or http://example.com/path?param=1 for the exact path /path and query parameter with the prefix param=1")
	spTTL := spCmd.Duration("ttl", time.Hour, "duration the signed request is valid")
	spFmt := spCmd.String("signature-format", "qp", "format to output.  Must be one of qp (to output query parameters to add to a URL), cookie (to output the cookie format), or pc (to output a full URL in path component format).")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "subcommand must be provided\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate-keys":
		genCmd.Parse(os.Args[2:])

		// Permission bits are 0600 since private keys should be private
		keyFile, err := os.OpenFile(*genKey, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open private key file for writing: %s\n", err)
			os.Exit(1)
		}
		pubFile, err := os.OpenFile(*genPub, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open public key file for writing: %s\n", err)
			os.Exit(1)
		}

		if err := genKeys(keyFile, pubFile); err != nil {
			fmt.Fprintf(os.Stderr, "could not generate keys: %s\n", err)
			os.Exit(1)
		}

	case "encode-key":
		ekCmd.Parse(os.Args[2:])

		pub, err := os.ReadFile(*ekPub)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read public key file: %s\n", err)
			os.Exit(1)
		}

		encodeKey(pub)

	case "sign-url":
		suCmd.Parse(os.Args[2:])

		if *suKeyset == "" {
			fmt.Fprintf(os.Stderr, "a keyset must be provided\n")
			os.Exit(1)
		}

		if *suURL == "" {
			fmt.Fprintf(os.Stderr, "a url must be provided\n")
			os.Exit(1)
		}

		key, err := os.ReadFile(*suKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read private key file: %s\n", err)
			os.Exit(1)
		}

		expiration := time.Now().Add(*suTTL)

		signURL(key, *suKeyset, *suURL, expiration)

	case "sign-prefix":
		spCmd.Parse(os.Args[2:])

		if *spKeyset == "" {
			fmt.Fprintf(os.Stderr, "a keyset must be provided\n")
			os.Exit(1)
		}

		if *spURL == "" {
			fmt.Fprintf(os.Stderr, "a url prefix must be provided\n")
			os.Exit(1)
		}

		key, err := os.ReadFile(*spKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read private key file: %s\n", err)
			os.Exit(1)
		}

		expiration := time.Now().Add(*spTTL)

		switch *spFmt {
		case "qp":
			signPrefixWithQueryParameters(key, *spKeyset, *spURL, expiration)
		case "cookie":
			signPrefixWithCookie(key, *spKeyset, *spURL, expiration)
		case "pc":
			// The path component mode only works if spURL is a path prefix
			// since we can't add path components after query parameters.
			if strings.ContainsRune(*spURL, '?') {
				fmt.Fprintf(os.Stderr, "the pc signature format does not work with query parameter prefixes\n")
				os.Exit(1)
			}

			signPrefixWithPathComponent(key, *spKeyset, *spURL, expiration)
		default:
			fmt.Fprintf(os.Stderr, "unknown signature-format: %q\n", *spFmt)
			os.Exit(1)
		}

	case "-h", "help", "-help", "--help":
		fmt.Fprintf(os.Stdout, "Usage: %s subcommand [subcommand args...]\n\nwhere: subcommand is one of generate-keys, encode-key, sign-url, sign-prefix, help\n", os.Args[0])
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
