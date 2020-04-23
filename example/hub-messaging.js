
const DELIMITER = '\n---\n';

function ParseMessage(callback) {
  let inputCache = '';
  return data => {
    inputCache += data.toString();
    let idx = inputCache.indexOf(DELIMITER);
    while (idx > -1) {
      const wireForm = inputCache.slice(0, idx);
      inputCache = inputCache.slice(idx + DELIMITER.length);
      const buff = Buffer.from(wireForm, 'base64');
      const msg = JSON.parse(buff.toString('ascii'));
      callback(msg);
      idx = inputCache.indexOf(DELIMITER);
    }
  };
}

function Envelope(msgObj) {
  const buff = Buffer.from(JSON.stringify(msgObj));
  return `${buff.toString('base64')}${DELIMITER}`;
}

function MessageWriter(outputStream) {
  let lastPromise = Promise.resolve();
  return message => {
    const nextPromise = new Promise(resolve => {
      const waitForDrain = !outputStream.write(Envelope(message));
      if (waitForDrain) outputStream.once('drain', resolve);
      else resolve();
    });
    lastPromise.then(nextPromise);
    lastPromise = nextPromise;
  };
}

module.exports = {
  DELIMITER,
  ParseMessage,
  Envelope,
  MessageWriter,
};
