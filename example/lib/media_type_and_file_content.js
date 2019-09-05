const readline = require('readline');
const { inform } = require('./log');
const eventsBuffer = {};

const reader = readline.createInterface({
  input: process.stdin
});

reader.on('line', (input) => {
  const event = JSON.parse(input);
  const acc = eventsBuffer[event.id] || {
    mediaType: null,
    hasMediaType: false,
    fileContent: null,
    hasFileContent: false,
  };

  if (event.role === 'media-type') {
    inform(event.id, 'media type received');
    acc.mediaType = event.payload;
    acc.hasMediaType = true;
  } else if (event.role === 'file-content') {
    inform(event.id, 'file content received');
    acc.fileContent = event.payload;
    acc.hasFileContent = true;
  }

  if (acc.hasMediaType && acc.hasFileContent) {
    inform(event.id, 'media type and file content received, so returning');
    const {status, body} = acc.fileContent;
    const {mediaType} = acc.mediaType;
    send(event.id, status, mediaType, body);
    delete eventsBuffer[event.id];
  } else {
    eventsBuffer[event.id] = acc;
  }
});

function send(id, status, contentType, body) {
  const event = {
    id,
    role: 'response',
    cmd: 'respond',
    payload: {
      status,
      contentType,
      body,
    },
  };
  const eventJson = JSON.stringify(event);
  process.stdout.write(`${eventJson}\n`, 'utf8');
}
