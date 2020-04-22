#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const {ParseMessage, MessageWriter} = require('./hub-messaging');
const write = MessageWriter(process.stdout);

// THIS IS NOT YET A STREAMING FILE READER

process.stdin.on('data', ParseMessage(({
  id,
  cmd, 
  url,
}) => {
  if (cmd === 'READ_FILE') {
    const filePath = (url === '') || (url === '/')
      ? '/index.html'
      : url;
    readFile(filePath, (err, data) => {
      if (err) {
        console.error('ERROR READING:', filePath);
        write({
          id,
          cmd: 'FAIL_WITH_NOT_FOUND',
          error: err.message,
        });
      } else {
        console.error('SUCCESS READING:', filePath, 'WITH DATA LENGTH:', data.length);
        write({
          id,
          cmd: 'PROCESS_FILE',
          encoding: 'base64',
          payload: data.toString('base64'),
        });
      }
    });
  }
}));

function readFile(url, callback) {
  // process.cwd()
  const fullPath = path.join('.', 'static', url);
  fs.readFile(fullPath, (err, data) => {
    if (err) callback(err);
    else callback(undefined, data);
  });
}