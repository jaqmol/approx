const readline = require('readline');
const fs = require('fs');
const path = require('path');
const {Header, Writer} = require('./utils');

const reader = readline.createInterface({
  input: process.stdin
});
const log = Writer({reader, stream: process.stderr});
const dispatch = Writer({reader, stream: process.stdout});

readMediaTypes((err, mediaTypeForExt) => {
  if (err) {
    log(`${err.toString()}\n`);
    process.exit(-1);
  } else {
    listen(mediaTypeForExt);
  }
});

function listen(mediaTypeForExt) {
  reader.on('line', (msgStr) => {
    const [header, body] = Header.parse(msgStr);
    if (!header || !body) return;
    
    const req = JSON.parse(body);
    const rawExt = path.extname(req.url.path);
    const extension = rawExt === '' ? '.html' : rawExt;
    const mediaType = mediaTypeForExt[extension];
    send(header.id, mediaType);
  });
}

function send(id, body) {
  const head = Header.stringify({id, role: 'media-type', mediaType: 'text/plain'});
  log(`Outbound message: ${head}${body}\n`);
  dispatch(`${head}${body}\n`);
}

function readMediaTypes(callback) {
  const csvFilePath = path.join(__dirname, 'media-types.csv');
  fs.readFile(csvFilePath, 'utf8', (err, data) => {
    if (err) {
      callback(err);
      return;
    }
    const acc = {};
    const lines = data.split('\n');
    for (let l of lines) {
      const cell = l.split(',');
      acc[cell[0]] = cell[1];
    }
    callback(null, acc);
  });
}
