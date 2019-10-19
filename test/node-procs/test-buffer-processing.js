#!/usr/bin/env node

const EventReader = require('../../node-approx/event-reader');
const ErrorWriter = require('../../node-approx/error-writer');
const EventWriter = require('../../node-approx/event-writer');

const input = new EventReader(process.stdin);
const error = new ErrorWriter(process.stderr);
const output = new EventWriter(process.stdout);

input.on('data', (chunk) => {
  const str = chunk.toString('utf8').replace(/\x00/g, '');
  
  try {
    const obj = JSON.parse(str);
    obj.first_name = obj.first_name.toUpperCase();
    obj.last_name = obj.last_name.toUpperCase();
    const lineStr = JSON.stringify(obj);
    const line = Buffer.from(lineStr);

    if (!output.write(line)) {
      input.pause();
      output.once('drain', () => input.resume());
    }
  } catch(err) {
    error.write(err);
  }
});

input.on('end', () => {
  output.end();
});

input.resume();


