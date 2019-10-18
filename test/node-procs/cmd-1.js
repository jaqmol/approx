#!/usr/bin/env node

const ApproxEventReader = require('./approx-event-reader');

const stdin = process.stdin,
      stdout = process.stdout,
      stderr = process.stderr;

const eventReader = new ApproxEventReader(stdin);

eventReader.on('data', (chunk) => {
  stderr.write(chunk);
  if (!stdout.write(chunk)) {
    // stdin.pause();
    eventReader.pause();
    stdout.once('drain', () => {
      // stdin.resume();
      eventReader.resume();
    });
  }
});

// stdin.on('data', (chunk) => {
//   if (!stdout.write(chunk)) {
//     stdin.pause();
//     stdout.once('drain', () => {
//       stdin.resume();
//     });
//   }
// });

// stdin.on('end', () => {
//   stdout.end();
// });
