const Event = {
  id: String,
	role: String,
	cmd: String,
	index: Number,
	sequence: Number,
  payload: Object,
}

const RequestPayload = {
	method: String,
	url: {
		host: String,
		path: String,
		query: {String: [String]},
	},
	headers: {String: [String]},
	body: String,
}