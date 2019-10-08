#!/usr/bin/env node

const stdin = process.stdin,
      stdout = process.stdout,
      stderr = process.stderr;

stdin.resume();

stdin.on('data', (chunk) => {
  // console.error(chunk);
  if (!stderr.write(chunk)) {
    stdin.pause();
    stderr.once('drain', () => {
      stdin.resume();
    });
  }
});

stdin.on('end', () => {
  stderr.end();
});
