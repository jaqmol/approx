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
    role: 'error',
    cmd: 'inform',
    payload: msg,
  };
  const infoJson = JSON.stringify(info, 2);
  process.stderr.write(`${infoJson}\n`, 'utf8');
}

function respond(msg) {
  const resp = {
    role: 'response',
    payload: {
      status: 200,
      contentType: 'application/json',
      body: msg
    },
  };
  const respJson = JSON.stringify(resp, 2);
  process.stdout.write(`${respJson}\n`, 'utf8');
}