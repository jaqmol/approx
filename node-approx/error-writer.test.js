const ErrorWriter = require('./error-writer');
const { Writable } = require('stream');
const { EvntEnd } = require('./config');

test('errors streaming', () => {
  const errorEvents = [];
  const originals = [];

  const output = new Writable({
    write(chunk, enc, cb) {
      errorEvents.push(chunk.toString('utf8'));
      cb();
    }
  });
  const errors = new ErrorWriter(output);
  
  for (const err of generateErrors(5000)) {
    originals.push(err);
    errors.write(err);
  }

  const received = errorEvents
    .map(e => e.slice(0, e.lastIndexOf(EvntEnd)))
    .map(s => JSON.parse(s));

  expect(received.length).toBe(originals.length);
  
  for (let i = 0; i < originals.length; i++) {
    const { message: originalMsg } = originals[i];
    const { message: receivedMsg } = received[i];
    expect(originalMsg).toEqual(receivedMsg);
  }
});

function* generateErrors(amount) {
  let index = 0;
  while (index < amount) {
    index++;
    try {
      throw new Error(`Test error: ${index + 1}`);
    } catch(e) {
      yield e;
    }
  }
}
