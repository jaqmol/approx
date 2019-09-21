import Message from './message';

test('Message from', () => {
  const id = '/';
  const role = 'test-message';
  const status = 200;
  const mediaType = 'text/plain';
  const encoding = 'utf8';
  const data = Buffer.from('TEST-DATA');

  const msg = Message.from({id, role, status, mediaType, encoding, data});
  const envelopeStr = msg.envelope().toString('utf8');
  expect(envelopeStr).toBe('51:/,test-message,0,true,200,text/plain,utf8;TEST-DATA');
  const envelope = Buffer.from(envelopeStr);
  const colonIdx = envelope.indexOf(':');
  const length = Number.parseInt(envelope.slice(0, colonIdx));
  const buffer = envelope.slice(colonIdx + 1);
  expect(length).toBe(buffer.length);
  
  const semicolonIdx = buffer.indexOf(';');
  const headerBuffer = buffer.slice(0, semicolonIdx);
  const headerComps = headerBuffer.toString('utf8').split(',');
  expect(headerComps[0]).toBe(id);
  expect(headerComps[1]).toBe(role);
  expect(headerComps[2]).toBe('0');
  expect(headerComps[3]).toBe('true');
  expect(headerComps[4]).toBe(String(status));
  expect(headerComps[5]).toBe(mediaType);
  expect(headerComps[6]).toBe(encoding);

  const dataBuffer = buffer.slice(semicolonIdx + 1);
  expect(data.equals(dataBuffer)).toBeTruthy();
});

test('Message parse', () => {
  const id = 'a1';
  const role = 'test-message';
  const status = 202;
  const mediaType = 'text/plain';
  const encoding = 'utf8';
  const data = Buffer.from('TEST-DATA');

  const origMsg = Message.from({id, role, status, mediaType, encoding, data});
  const envelope = origMsg.envelope();
  const msgBuff = envelope.slice(envelope.indexOf(':') + 1);
  
  const msg = Message.parse({buffer: msgBuff});
  expect(msg.id).toBe(id);
  expect(msg.role).toBe(role);
  expect(msg.status).toBe(status);
  expect(msg.mediaType).toBe(mediaType);
  expect(msg.encoding).toBe(encoding);
  expect(msg.data.equals(data)).toBeTruthy();
});
