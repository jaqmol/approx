#!/usr/bin/env node

const http = require('http');
const path = require('path');
const cuid = require('cuid');
const {SendMessage, ParseMessage} = require('./hub-messaging');
const send = SendMessage(process.stdout);

const port = 3000;
const pending = {};

const server = http.createServer((request, response) => {
  const id = cuid();
  
  pending[id] = {id, request, response, needsHeaders: true};

  response.setTimeout(10000);

  send({id, cmd: 'READ_FILE', url: request.url});
  send({id, cmd: 'FIND_MEDIA_TYPE', ext: path.extname(request.url)});
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
  if (cmd === 'PROCESS_FILE_CHUNK') {
    const {response, needsHeaders} = pending[id];
    if (needsHeaders) {
      response.setHeader('Content-Type', contentType);
      pending[id].needsHeaders = false;
    }
    const data = Buffer.from(payload, encoding);
    response.write(data);
  } else if (cmd === 'CONCLUDE_FILE') {
    const {response} = pending[id];
    response.end();
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
