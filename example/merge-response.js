#!/usr/bin/env node

const {ParseMessage, SendMessage} = require('./hub-messaging');
const send = SendMessage(process.stdout);

const contentTypeForID = {};
const chunkMsgsForID = {};

process.stdin.on('data', ParseMessage(({
  id,
  cmd, 
  mediaType,
  encoding,
  payload,
  error,
}) => {
  if (cmd === 'FAIL_WITH_NOT_FOUND') {
    send({id, cmd, error});
    delete contentTypeForID[id];
    delete chunkMsgsForID[id];
  } else if (cmd === 'PROCESS_MEDIA_TYPE') {
    contentTypeForID[id] = mediaType;
    flushIfPossible(id);
  } else if (cmd === 'PROCESS_FILE_CHUNK') {
    addChunkMsg(id, {id, cmd, encoding, payload});
    flushIfPossible(id);
  } else if (cmd === 'CONCLUDE_FILE') {
    addChunkMsg(id, {id, cmd});
    flushIfPossible(id);
    delete contentTypeForID[id];
    delete chunkMsgsForID[id];
  }
}));

function addChunkMsg(id, msg) {
  const msgs = chunkMsgsForID[id] || [];
  msgs.push(msg);
  chunkMsgsForID[id] = msgs;
}

function flushIfPossible(id) {
  const contentType = contentTypeForID[id];
  const msgs = chunkMsgsForID[id] || [];
  if ( (typeof contentType !== 'undefined') && 
       (typeof msgs !== 'undefined') )
  {
    msgs.map(msg => ({...msg, contentType}))
        .forEach(send);
    chunkMsgsForID[id] = [];
  }
}
