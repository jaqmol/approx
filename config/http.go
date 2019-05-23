package config

// HTTP ...
type HTTP struct {
	Endpoint string
	Proxy    Proxy
}

// HTTPS ...
type HTTPS struct {
	Endpoint    string
	Proxy       Proxy
	Certificate Certificate
}

// Proxy ...
type Proxy struct {
	Out string
	In  string
}

/*
"staticFileServer": {
    "type": "http",
    "endpoint": {"config": "endpoint"},
    "proxy": {
      "out": "_copyRequest",
      "in": "_bufferUntilMimeType"
    }
	},
*/
