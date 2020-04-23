#!/usr/bin/env node

const http = require('http');
const path = require('path');
const cuid = require('cuid');
const {MessageWriter, ParseMessage} = require('./hub-messaging');
const write = MessageWriter(process.stdout);

const port = 3000;
const pending = {};

const server = http.createServer((request, response) => {
  const id = cuid();
  
  pending[id] = {id, request, response};

  write({id, cmd: 'READ_FILE', url: request.url});
  write({id, cmd: 'FIND_MEDIA_TYPE', ext: path.extname(request.url)});
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
  error,
}) => {
  if (cmd === 'RESPOND') {
    const {response} = pending[id];
    const data = Buffer.from(payload, encoding);
    response.setHeader('Content-Type', contentType);
    response.end(data);
    delete pending[id];
  } else if (cmd === 'FAIL_WITH_NOT_FOUND') {
    const {response} = pending[id];
    response.writeHead(404, { 'Content-Type': 'text/html' });
    response.end(`<h1>NOT FOUND</h1><p>${error}</p>`, 'utf-8')
    delete pending[id];
  } else {
    console.error('Web server: Unknown message:', cmd);
  }
}));
