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

function MakeMessage(msgObj) {
  const buff = Buffer.from(JSON.stringify(msgObj));
  return `${buff.toString('base64')}${DELIMITER}`;
}

module.exports = {
  DELIMITER,
  ParseMessage,
  MakeMessage,
};
