const readline = require('readline');
const {Header, Writer} = require('./utils');
const buff = {};

const reader = readline.createInterface({
  input: process.stdin
});

const log = Writer({reader, stream: process.stderr});
const dispatch = Writer({reader, stream: process.stdout});

reader.on('line', (msgStr) => {
  const [header, body] = Header.parse(msgStr);
  if (!header || !body) return;

  const acc = buff[header.id] || {
    mediaType: null,
    hasMediaType: false,
    body: null,
    status: 0,
    hasBodyAndStatus: false,
  };

  if (header.role === 'media-type') {
    acc.mediaType = body;
    acc.hasMediaType = true;
    // log(`did receive media type: ${body}\n`);
  } else if (header.role === 'file-content') {
    acc.body = body;
    acc.status = header.status;
    acc.hasBodyAndStatus = true;
    // log(`did receive file content (${header.status}): ${body}\n`);
  }

  if (acc.hasMediaType && acc.hasBodyAndStatus) {
    send(header.id, acc.status, acc.mediaType, acc.body);
    delete buff[header.id];
  } else {
    buff[header.id] = acc;
  }
});

function send(id, status, mediaType, body) {
  const head = Header.stringify({id, role: 'response', status, mediaType, encoding: 'base64'});
  // log(`About to write response: ${head}${body}\n`);
  dispatch(`${head}${body}\n`);
}
