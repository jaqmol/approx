#!/usr/bin/env node

const {ParseMessage, MessageWriter} = require('./hub-messaging');
const write = MessageWriter(process.stdout);
const pending = {};

process.stdin.on('data', ParseMessage(({
  id,
  cmd, 
  mediaType,
  encoding,
  payload,
  error,
}) => {
  console.error('TRYING TO MERGE:', cmd);
  if (cmd === 'FAIL_WITH_NOT_FOUND') {
    respondWithNotFound(cmd, id, error);
  } else {
    const resp = pending[id] || {id, cmd: 'RESPOND'};
    if (cmd === 'PROCESS_MEDIA_TYPE') {
      resp.contentType = mediaType;
    } else if (cmd === 'PROCESS_FILE') {
      resp.encoding = encoding;
      resp.payload = payload;
    }
    pending[id] = resp;
    respondIfComplete(id);
  }
}));

function respondIfComplete(id) {
  const resp = pending[id];
  const hasContentType = typeof resp.contentType !== 'undefined';
  if (hasContentType && resp.encoding && resp.payload) {
    console.error('MERGE COMPLETE:', resp.id);
    write(resp);
    delete pending[id];
  } else {
    console.error('MERGE INCOMPLETE:', resp.id);
  }
}

function respondWithNotFound(cmd, id, error) {
  console.error('MERGE PASSING ON NOT FOUND:', id);
  write({id, cmd, error});
  delete pending[id];
}
