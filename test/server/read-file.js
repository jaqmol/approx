const fs = require('fs');
const path = require('path');
const readline = require('readline');
const {Header, Writer} = require('./utils');

const reader = readline.createInterface({
  input: process.stdin
});

const log = Writer({reader, stream: process.stderr});
const dispatch = Writer({reader, stream: process.stdout});

let dataSentCount = 0;

reader.on('line', (msgStr) => {
  const [header, body] = Header.parse(msgStr);
  if (!header || !body) return;
  // log(`Inbound header: ${JSON.stringify(header)}\n`);
    
  const req = JSON.parse(body);
  let filePath = req.url.path;
  if (filePath === '/') {
    filePath = '/index.html';
  }

  readFileB64(filePath, (err, body) => {
    if (err) {
      if (err.code === 'ENOENT') {
        send(header.id, 404, '', '');
        // log(`File not found @ ${filePath}\n`);
      } else {
        const errStr = err.toString();
        send(header.id, 500, 'utf8', errStr);
        // log(`Error reading file @ ${filePath}: ${errStr}\n`);
      }
    } else {
      send(header.id, 200, 'base64', body);
      // log(`Did read file @ ${filePath}\n`);
    }
  });
});

function send(id, status, encoding, body) {
  const head = Header.stringify({id, role: 'file-content', status, encoding});
  dispatch(`${head}${body}\n`);
  dataSentCount++;
  log(`${dataSentCount}: Did send data of ${id}\n`);
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