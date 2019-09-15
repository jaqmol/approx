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

  let acc;
  if (buff.has(header.id)) {
    acc = buff.get(header.id);
  } else {
    acc = {
      mediaType: null,
      chunks: [],
    };
  }

  if (header.role === 'media-type') {
    mediaTypeCount++;
    acc.mediaType = body;
  } else if (header.role === 'file-content') {
    bodyAndStatusCount++;
    const {seq, status, isEnd} = header;
    acc.chunks.push({seq, status, isEnd, body});
  }

  log(`mediaTypeCount: ${mediaTypeCount}, bodyAndStatusCount: ${bodyAndStatusCount}\n`);

  if (acc.mediaType && acc.chunks.length) {
    log(`combine responds with ${acc.chunks.length} chunks\n`);
    acc.chunks = acc.chunks.filter(({seq, status, body, isEnd}) => {
      send(header.id, seq, status, acc.mediaType, body, isEnd);
      return false;
    });
  }

  if (header.isEnd && (header.role === 'file-content')) {
    log(`combine deletes accumulator for id ${header.id}\n`);
    buff.delete(header.id);
  } else {
    log(`combine updates accumulator for id ${header.id}\n`);
    buff.set(header.id, acc);
  }
});

function send(id, seq, status, mediaType, body, isEnd) {
  const head = Header.stringify({id, seq, role: 'response-chunk', status, mediaType, encoding: 'base64', isEnd});
  dispatch(`${head}${body}\n`);
}
