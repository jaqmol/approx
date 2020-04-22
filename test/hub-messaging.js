let inputCache = '';

const DELIMITER = '\n---\n';

function ParseMessage(callback) {
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
  const jobsCache = [];
  let running = false;
  const next = () => {
    if (running) return;
    if (!jobsCache.length) return;
    const perform = jobsCache.splice(0, 1)[0];
    running = true;
    perform();
  };
  const conclude = () => {
    running = false;
    next();
  };
  return message => {
    jobsCache.push(() => {
      const waitForDrain = !outputStream.write(Envelope(message));
      if (waitForDrain) outputStream.once('drain', conclude);
      else setTimeout(conclude, 0);
    });
    next();
  };
}

module.exports = {
  DELIMITER,
  ParseMessage,
  Envelope,
  MessageWriter,
};
