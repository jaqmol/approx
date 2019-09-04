const readline = require('readline');

const reader = readline.createInterface({
  input: process.stdin
});

const collector = {};

reader.on('line', (input) => {
  const msg = JSON.parse(input);
  let allMessages = collector[msg.id];
  if (!allMessages) {
    allMessages = [];
  }
  allMessages.push(msg);
  if (allMessages.length === 2) {
    respond(msg.id, allMessages);  
    delete collector[msg.id];
  } else {
    collector[msg.id] = allMessages;
  }

  // inform(msg);
  // respond(msg);
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

function respond(id, allMsgs) {
  let bodyStr = JSON.stringify(allMsgs, null, 2);
  let body = Buffer.from(bodyStr).toString('base64');
  const resp = {
    id,
    role: 'response',
    payload: {
      status: 200,
      contentType: 'application/json',
      body,
    },
  };
  const respJson = JSON.stringify(resp, 2);
  writeTo(process.stdout, respJson);
}

const writeTo = (stream, message) => stream.write(`${message}\n`, 'utf8');