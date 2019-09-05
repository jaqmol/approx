const fs = require('fs');
const path = require('path');
const readline = require('readline');
const { inform, exit } = require('./log');

readMediaTypes((err, mediaTypeForExt) => {
  if (err) {
    exit(null, err.toString());
  } else {
    listen(mediaTypeForExt);
  }
});

function listen(mediaTypeForExt) {
  const reader = readline.createInterface({
    input: process.stdin
  });

  reader.on('line', (input) => {
    const event = JSON.parse(input);
    const rawExt = path.extname(event.payload.url.path);
    const extension = rawExt === '' ? '.html' : rawExt;
    inform(event.id, "checking media type of extension " + extension);
    const mediaType = mediaTypeForExt[extension];
    send(event.id, event.index, extension, mediaType);
  });
}

function send(id, index, extension, mediaType) {
  const event = {
    id,
    index,
    role: 'media-type',
    cmd: 'set-content-type',
    payload: {
      extension,
      mediaType,
    },
  };
  const eventJson = JSON.stringify(event);
  process.stdout.write(`${eventJson}\n`, 'utf8');
}

function readMediaTypes(callback) {
  const csvFilePath = path.join(__dirname, 'media-type-for-extension.csv');
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