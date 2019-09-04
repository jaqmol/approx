const readline = require('readline');

const reader = readline.createInterface({
  input: process.stdin
});

reader.on('line', (input) => {
  const msg = JSON.parse(input);
  // inform(msg);
  respond(msg);
});

function inform(msg) {
  const info = {
    id: msg.id,
    role: 'log',
    cmd: 'inform',
    payload: "\"fork-merge-and-return.js got an event\"",
  };
  const infoJson = JSON.stringify(info, 2);
  writeTo(process.stderr, infoJson);
}

function respond(msg) {
  let originalRequestJson = JSON.stringify(msg, null, 2);
  let originalRequest = Buffer.from(originalRequestJson).toString('base64');
  const resp = {
    id: msg.id,
    index: msg.index,
    role: 'passthrought',
    payload: {
      originalRequest,
      message: `Hello from passthrough ${msg.index === 0 ? 'A' : 'B'}`,
    },
  };
  const respJson = JSON.stringify(resp, 2);
  writeTo(process.stdout, respJson);
}

const writeTo = (stream, message) => stream.write(`${message}\n`, 'utf8');
