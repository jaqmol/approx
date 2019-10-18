#!/usr/bin/env node

const stdin = process.stdin,
      stdout = process.stdout,
      stderr = process.stderr;

stdin.on('data', (chunk) => {
  console.log(chunk.toString('utf8'));
  // if (!stderr.write(chunk)) {
  //   stdin.pause();
  //   stderr.once('drain', () => {
  //     stdin.resume();
  //   });
  // }
});

stdin.on('end', () => {
  stderr.end();
});

stdin.resume();


