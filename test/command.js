#!/usr/bin/env node

const stdin = process.stdin,
      stdout = process.stdout;

stdin.resume();

stdin.on('data', (chunk) => {
  console.error(chunk);
  if (!stdout.write(chunk)) {
    stdin.pause();
    stdout.once('drain', () => {
      stdin.resume();
    });
  }
});

stdin.on('end', () => {
  stdout.end();
});
