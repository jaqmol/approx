#!/usr/bin/env node

const {ParseMessage, MakeMessage} = require('./hub-messaging');
const pending = {};

process.stdin.on('data', ParseMessage(({
  id,
  cmd, 
  mediaType,
  encoding,
  payload,
  error,
}) => {
  const resp = pending[id] || {id, cmd: 'RESPOND'};
  if (cmd === 'PROCESS_MEDIA_TYPE') {
    resp.contentType = mediaType;
  } else if (cmd === 'PROCESS_FILE') {
    resp.encoding = encoding;
    resp.payload = payload;
  } else if (cmd === 'PROCESS_FILE_ERROR') {
    resp.cmd = 'FAIL';
    resp.error = error;
  }
  pending[id] = resp;
  respondIfComplete(id);
}));

function respondIfComplete(id) {
  const resp = pending[id];
  if (resp.cmd === 'FAIL') {
    process.stdout.write(MakeMessage(resp));
    delete pending[id];
  } else {
    const hastMediaType = typeof resp.mediaType !== 'undefined';
    if (hastMediaType && resp.encoding && resp.payload) {
      process.stdout.write(MakeMessage(resp));
      delete pending[id];
    }
  }
}
