const readline = require('readline');
const {Header, Writer} = require('./utils');
const buff = new Map();

const reader = readline.createInterface({
  input: process.stdin
});

const log = Writer({reader, stream: process.stderr});
const dispatch = Writer({reader, stream: process.stdout});

let mediaTypeCount = 0;
let bodyAndStatusCount = 0;

reader.on('line', (msgStr) => {
  const [header, body] = Header.parse(msgStr);
  if (!header || !body) return;

  // log(`header.id: ${header.id}\n`);

  let acc;
  if (buff.has(header.id)) {
    acc = buff.get(header.id);
  } else {
    acc = {
      mediaType: null,
      hasMediaType: false,
      body: null,
      status: 0,
      hasBodyAndStatus: false,
    };
  }

  if (header.role === 'media-type') {
    mediaTypeCount++;
    acc.mediaType = body;
    acc.hasMediaType = true;
  } else if (header.role === 'file-content') {
    bodyAndStatusCount++;
    acc.body = body;
    acc.status = header.status;
    acc.hasBodyAndStatus = true;
  }

  log(`mediaTypeCount: ${mediaTypeCount}, bodyAndStatusCount: ${bodyAndStatusCount}\n`);
  
  if (acc.hasMediaType && acc.hasBodyAndStatus) {
    send(header.id, acc.status, acc.mediaType, acc.body);
    buff.delete(header.id);
  } else {
    buff.set(header.id, acc);
  }

  // for (var [key, value] of buff.entries()) {
  //   log(`${key}: hasMediaType: ${value.hasMediaType}, hasBodyAndStatus: ${value.hasBodyAndStatus}\n`);
  // }
});

function send(id, status, mediaType, body) {
  const head = Header.stringify({id, role: 'response', status, mediaType, encoding: 'base64'});
  dispatch(`${head}${body}\n`);
}
