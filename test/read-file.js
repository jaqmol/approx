#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const {ParseMessage, MakeMessage} = require('./hub-messaging');

// THIS IS NOT YET A STREAMING FILE READER

process.stdin.on('data', ParseMessage(({
  id,
  cmd, 
  url,
}) => {
  if (cmd === 'READ_FILE') {
    readFile(url, (err, data) => {
      if (err) {
        process.stdout.write(MakeMessage({
          id,
          cmd: 'PROCESS_FILE_ERROR',
          error: err.message,
        }));
      } else {
        console.error('DID READ:', url, 'WITH DATA LENGTH:', data.length);
        process.stdout.write(MakeMessage({
          id,
          cmd: 'PROCESS_FILE',
          encoding: 'base64',
          payload: data.toString('base64'),
        }));
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