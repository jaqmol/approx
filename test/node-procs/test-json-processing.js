#!/usr/bin/env node

const JSONEventReader = require('../../node-approx/json-event-reader');
const LogWriter = require('../../node-approx/log-writer');
const JSONEventWriter = require('../../node-approx/json-event-writer');

const input = new JSONEventReader(process.stdin);
const log = new LogWriter(process.stderr);
const output = new JSONEventWriter(process.stdout);

input.on('data', (obj) => {
  try {
    obj.first_name = obj.first_name.toUpperCase();
    obj.last_name = obj.last_name.toUpperCase();
    log.info('Did process:', obj.id);

    if (!output.write(obj)) {
        input.pause();
      output.once('drain', () => input.resume());
    }
  } catch(err) {
    log.error(err);
  }
});

input.on('end', () => {
  output.end();
});

input.resume();


