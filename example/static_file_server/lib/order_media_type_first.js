const axmsg = require('./axmsg');
const axenvs = require('./axenvs');
const errors = new axmsg.Errors('order_media_type_first');
const envs = new axenvs.Envs('order_media_type_first');

const [ins, outs] = envs.insOuts();
let fatalErrors = false;
if (ins.length !== 1) {
  errors.log(null, new Error(`Exactly 1 input expected, but got: ${ins.length}`));
  fatalErrors = true;
}
if (outs.length !== 1) {
  errors.log(null, new Error(`Exactly 1 output expected, but got: ${out.length}`));
  fatalErrors = true;
}
if (fatalErrors) {
  errors.logFatal(null, new Error("Too many fatal errors, can't run"));
}
startListening(ins[0], outs[0]);

function startListening(reader, writer) {
  const actionsBuffer = composeBuffer(writer);
  reader.on(action => {
    // input example:
    // {"axmsg":1,"id":21,"role":"media-type","cmd":"set-content-type","data":{"extension":".html","mediaType":"text/html"}}
    // {"axmsg":1,"id":21,"role":"response-data","cmd":"respond-with-data","data":{"status":200,"bodyB64":"bG9vaywgdGhlcmUncyBubyBvdXRzaWRl"}}
    // {"axmsg":1,"id":21,"role":"response-data","cmd":"respond-with-error","data":{"status":404}}
    if (actionsBuffer.addMediaTypeActionIfViable(action)) {
      actionsBuffer.flushIfViable();
    } else if (actionsBuffer.addResponseDataActionIfViable(action)) {
      actionsBuffer.flushIfViable();
    } else {
      if (action.id) {
        actionsBuffer.clearIDs([action.id]);
      }
      errors.log(action.id, new Error(`Unknown action: ${JSON.stringify(action)}`));
    }
  });
}

function composeBuffer(writer) {
  const mediaTypeActions = {};
  const responseDataActions = {};

  const isMediaTypeAction = action => {
    return action.axmsg === 1 && action.role === 'media-type' && action.cmd === 'set-content-type';
  };
  
  const isResponseDataAction = action => {
    return action.axmsg === 1 && action.role === 'response-data';
  };

  const addMediaTypeActionIfViable = action => {
    if (!isMediaTypeAction(action)) return false;
    mediaTypeActions[action.id] = action;
    return true;
  };

  const addResponseDataActionIfViable = action => {
    if (!isResponseDataAction(action)) return false;
    responseDataActions[action.id] = action;
    return true;
  };

  const flushIfViable = () => {
    const flushableIDs = Object.keys(mediaTypeActions);
    const clearableIDs = [];
    for (let id of flushableIDs) {
      const responseData = responseDataActions[id];
      if (responseData) {
        const mediaType = mediaTypeActions[id];
        writer.writeAction(mediaType);
        writer.writeAction(responseData);
        clearableIDs.push(id);
      }
    }
    clearIDs(clearableIDs);
  };

  const clearIDs = ids => {
    for (let id of ids) {
      delete mediaTypeActions[id];
      delete responseDataActions[id];
    }
  };

  return {
    addMediaTypeActionIfViable,
    addResponseDataActionIfViable,
    flushIfViable,
    clearIDs,
  };
}