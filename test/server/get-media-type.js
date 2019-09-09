const readline = require('readline');
const fs = require('fs');
const path = require('path');
const writer = require('./writer');

const reader = readline.createInterface({
  input: process.stdin
});
const log = writer(reader, process.stderr);
const dispatch = writer(reader, process.stdout);

reader.on('line', (msgStr) => {
  const [header, body] = parseHeader(msgStr);
  if (!header || !body) return;
  log(msgStr);
  respond(header.id, body);
});


function respond(id, msgStr) {
  const head = stringifyHeader(id, 'response', 0, true, 200, 'application/json', '');
  dispatch(`${head}${msgStr}\n`);
}

function parseHeader(msgStr) {
  if (!msgStr) return [];
  const semicolonIdx = msgStr.indexOf(';');
  const headerStr = msgStr.slice(0, semicolonIdx);
  const comps = headerStr.split(',');
  return [
    {
      id: comps[0],
      role: comps[1],
      isEnd: JSON.parse(comps[2]),
      mediaType: comps[3],
      encoding: comps[4],
    },
    msgStr.slice(semicolonIdx + 1),
  ];
}

const stringifyHeader = (id, role, seq, isEnd, status, mediaType, encoding) => 
  `${id},${role},${seq},${isEnd},${status},${mediaType},${encoding};`;

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
