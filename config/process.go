package config

// Process ...
type Process struct {
	Path        string
	In          string
	Out         string
	Environment map[string]string
}

/*
"_bufferUntilMimeType": {
    "type": "process",
    "path": "lib/buffer-until",
    "in": "_zipResponse",
    "environment": {
      "PROP_PATH": "equals(header.result.cmd,set-mime-type)"
    }
	},
*/
