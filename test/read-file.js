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
    readFileB64(url, (err, data) => {
      if (err) {
        process.stdout.write(MakeMessage({
          id,
          cmd: 'PROCESS_FILE_ERROR',
          error: err.message,
        }));
      } else {
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

function readFileB64(filePath, callback) {
  const fullPath = path.join(process.cwd(), process.env.BASE_PATH, filePath);
  fs.readFile(fullPath, (err, data) => {
    if (err) callback(err);
    else callback(undefined, data);
  });
}