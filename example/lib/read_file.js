const fs = require('fs');
const path = require('path');
const readline = require('readline');
const { inform, exit } = require('./log');

const reader = readline.createInterface({
  input: process.stdin
});

reader.on('line', (input) => {
  const event = JSON.parse(input);
  let filePath = event.payload.url.path;
  if (filePath === '/') {
    filePath = '/index.html';
  }
  // inform(event.id, `reading file ${filePath}`);
  readFileB64(filePath, (err, body) => {
    if (err) {
      if (err.code === 'ENOENT') {
        send(event.id, event.index, 404);
      } else {
        exit(event.id, err.toString());
        send(event.id, event.index, 500, err.toString());
      }
    } else {
      // inform(event.id, `did read file ${filePath}, body: ${body}`);
      send(event.id, event.index, 200, body);
    }
  });
});

function send(id, index, status, body) {
  const event = {
    id,
    index,
    role: 'file-content',
    cmd: 'respond',
    payload: {
      status,
      body,
    },
  };
  const eventJson = JSON.stringify(event);
  process.stdout.write(eventJson, 'utf8');
  process.stdout.write('\n', 'utf8');
}

function readFileB64(filePath, callback) {
  const fullPath = path.join(process.cwd(), process.env.BASE_PATH, filePath);
  fs.readFile(fullPath, (err, data) => {
    if (err) {
      callback(err);
      return;
    }
    const buf = Buffer.from(data);
    callback(null, buf.toString('base64'));
  });
}