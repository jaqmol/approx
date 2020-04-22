#!/usr/bin/env node

const http = require('http');
const path = require('path');
const cuid = require('cuid');
const {MakeMessage, ParseMessage} = require('./hub-messaging');

const port = 3000;
const pending = {};

const server = http.createServer((request, response) => {
  const id = cuid();
  
  pending[id] = {id, request, response};

  process.stdout.write(MakeMessage({id, cmd: 'READ_FILE', url: request.url}));
  process.stdout.write(MakeMessage({id, cmd: 'FIND_MEDIA_TYPE', ext: path.extname(request.url)}));
});

server.listen(port, (err) => {
  if (err) {
    return console.error('something bad happened', err);
  }
});

process.stdin.on('data', ParseMessage(({
  id, 
  cmd, 
  encoding,
  payload, 
  contentType,
}) => {
  if (cmd === 'RESPOND') {
    const {response} = pending[id];
    const data = Buffer.from(payload, encoding);
    response.setHeader('Content-Type', contentType);
    response.end(data);
    delete pending[id];
  } else {
    console.error('Web server: Unknown message:', cmd);
  }
}));
