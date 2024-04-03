package main

var PUBLIC_URIS = []string{
	"/crate-api/subscriber/sign-up",
	"/crate-api/subscriber/sign-in",
	"/crate-api/subscriber/validate-token",
	// `^/crate-api/subscriber/[a-zA-Z0-9-]+/[0-9]+$`,
}

const HEADER_API_VERSION = "x-api-version"
