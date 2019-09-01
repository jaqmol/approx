const readline = require('readline');

const reader = readline.createInterface({
  input: process.stdin
});

reader.on('line', (input) => {
  const msg = JSON.parse(input);
  inform(msg);
  respond(msg);
});

function inform(msg) {
  const info = {
    id: msg.id,
    role: 'error',
    cmd: 'inform',
    payload: "Hi from node",
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