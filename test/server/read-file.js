const fs = require('fs');
const path = require('path');
const readline = require('readline');
const {Header, Writer} = require('./utils');

const reader = readline.createInterface({
  input: process.stdin
});

const dispatch = Writer({reader, stream: process.stdout});

reader.on('line', (msgStr) => {
  const [header, body] = Header.parse(msgStr);
  if (!header || !body) return;
    
  const req = JSON.parse(body);
  let filePath = req.url.path;
  if (filePath === '/') {
    filePath = '/index.html';
  }

  let seq = 0;
  readFileB64(filePath, (err, data, isEnd) => {
    if (err) {
      if (err.code === 'ENOENT') {
        send(header.id, -1, 404, '', '', isEnd);
      } else {
        const errStr = err.toString();
        send(header.id, -1, 500, 'utf8', errStr, isEnd);
      }
    } else {
      send(header.id, seq, 200, 'base64', data, isEnd);
    }
    seq++;
  });
});

function send(id, seq, status, encoding, body, isEnd) {
  const head = Header.stringify({id, seq, role: 'file-content', status, encoding, isEnd});
  dispatch(`${head}${body}\n`);
}

function readFileB64(filePath, callback) {
  const fullPath = path.join(process.cwd(), process.env.BASE_PATH, filePath);
  const reader = fs.createReadStream(fullPath);
  let lastData = null;
  reader
    .on('error', (err) => {
      callback(err, null, true);
    }) 
    .on('data', (chunk) => {
      if (lastData) {
        callback(null, lastData, false);
      }
      lastData = chunk.toString('base64');
    })
    .on('end', () => {
      callback(null, lastData, true);
    });
}