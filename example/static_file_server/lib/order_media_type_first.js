const fs = require('fs');
const axmsg = require('./axmsg');
// const errors = new axmsg.Errors('read_file');
const reader = new axmsg.Reader(process.stdin);
const writer = new axmsg.Writer(process.stdout);

let mediaTypeForId = {};
let dataForId = {};

reader.on((action) => {
  // input example:
  // {"axmsg":1,"id":21,"role":"media-type","cmd":"set-content-type","data":{"extension":".html","mediaType":"text/html"}}
  // {"axmsg":1,"id":21,"role":"response-data","cmd":"respond-with-data","data":{"status":200,"bodyB64":"bG9vaywgdGhlcmUncyBubyBvdXRzaWRl"}}
  // {"axmsg":1,"id":21,"role":"response-data","cmd":"respond-with-error","data":{"status":404}}
  if (action.axmsg === 1) {
    if (action.role === 'media-type' && action.cmd === 'set-content-type') {
      mediaTypeForId[action.id] = action;
    } else if (action.role === 'response-data') {
      dataForId[action.id] = action;
    }
  }
  const flushedIDs = flushIfPossible();
  clearAllForIDs(flushedIDs);
});

function flushIfPossible() {
  const acc = [];
  for (var id in mediaTypeForId) {
    const data = dataForId[id];
    if (data) {
      if (data.cmd === 'respond-with-data') {
        const mediaType = mediaTypeForId[id];
        writer.writeAction(mediaType  );
      }
      writer.writeAction(data);
      acc.push(id);
    }
  }
  return acc;
}

function clearAllForIDs(ids) {
  for (let id of ids) {
    mediaTypeForId[id] = undefined;
    dataForId[id] = undefined;
  }
}