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
  process.stderr.write(`${infoJson}\n`, 'utf8');
}

function respond(msg) {
  let bodyStr = JSON.stringify(msg, null, 2);
  let body = Buffer.from(bodyStr).toString('base64');
  const resp = {
    id: msg.id,
    role: 'response',
    payload: {
      status: 200,
      contentType: 'application/json',
      body,
    },
  };
  const respJson = JSON.stringify(resp, 2);
  process.stdout.write(`${respJson}\n`, 'utf8');
}