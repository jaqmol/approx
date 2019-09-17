const {Message, ChunkCache} = require('./message.js');

test('Message from', () => {
  const id = 'a1';
  const role = 'test-message';
  const status = 200;
  const mediaType = 'text/plain';
  const encoding = 'utf8';
  const data = Buffer.from('TEST-DATA');

  const msg = Message.from({id, role, status, mediaType, encoding, data});
  const envelopeStr = msg.toBuffer().toString('utf8');
  expect(envelopeStr).toBe('53:a1,test-message,0,true,200,text/plain,utf8,;TEST-DATA');
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
  const envelope = origMsg.toBuffer();
  const msgBuff = envelope.slice(envelope.indexOf(':') + 1);
  
  const msg = Message.parse({buffer: msgBuff});
  expect(msg.header.id).toBe(id);
  expect(msg.header.role).toBe(role);
  expect(msg.header.status).toBe(status);
  expect(msg.header.mediaType).toBe(mediaType);
  expect(msg.header.encoding).toBe(encoding);
  expect(msg.data.equals(data)).toBeTruthy();
});

test('ChunkCache simple test', () => {
  const data = Buffer.from('TEST-DATA');
  const lengthBuff = Buffer.from(`${data.length}:`);
  const envelope = Buffer.concat([lengthBuff, data]);
  const cache = new ChunkCache();
  cache.push(envelope);
  const msgBuff = cache.nextMessageBuffer();
  expect(data.equals(msgBuff)).toBeTruthy();
});

const letterNames = 'ALPHA,BETA,GAMMA,DELTA,EPSILON,ZETA,ETA,THETA,IOTA,KAPPA,LAMBDA,MU,NU,XI,OMICRON,PI,RHO,SIGMA,TAU'.split(',');

function chunkCacheTestEnvelope(content) {
  const data = Buffer.from(content);
  const lengthBuff = Buffer.from(`${data.length}:`);
  return [data, Buffer.concat([lengthBuff, data])];
}

test('ChunkCache double chunk test', () => {
  const [dataA, envelopeA] = chunkCacheTestEnvelope(`TEST-DATA-${letterNames[0]}`);
  const [dataB, envelopeB] = chunkCacheTestEnvelope(`TEST-DATA-${letterNames[1]}`);

  const cache = new ChunkCache();
  cache.push(Buffer.concat([envelopeA, envelopeB]));
  
  const msgBuffA = cache.nextMessageBuffer();
  expect(dataA.equals(msgBuffA)).toBeTruthy();

  const msgBuffB = cache.nextMessageBuffer();
  expect(dataB.equals(msgBuffB)).toBeTruthy();
});

function chunkCacheTestChopEnvelope(buffer, count) {
  const acc = [];
  const pieceLen = buffer.length / count;
  for (let i = 0; i < count; i++) {
    acc.push(buffer.slice(i * pieceLen, (i + 1) * pieceLen));
  }
  return acc;
}

test('ChunkCache pieced chunks test', () => {
  const datasAndEnvelopes = letterNames.map((ln) => {
    return chunkCacheTestEnvelope(`TEST-DATA-${ln}`);
  });
  
  const originDatas = datasAndEnvelopes.map(([data]) => data);
  const envelopes = datasAndEnvelopes.map(([_, env]) => env);

  const combinedEnvelopes = Buffer.concat(envelopes);
  const chunks = chunkCacheTestChopEnvelope(
    combinedEnvelopes, 
    envelopes.length + (envelopes.length / 5),
  );

  const cache = new ChunkCache();
  for (const c of chunks) {
    cache.push(c);
  }

  for (const od of originDatas) {
    const msgBuff = cache.nextMessageBuffer();
    expect(od.equals(msgBuff)).toBeTruthy();
  }
});

test('ChunkCache rattled chunks test', () => {
  const datasAndEnvelopes = letterNames.map((ln) => {
    return chunkCacheTestEnvelope(`TEST-DATA-${ln}`);
  });
  
  const originDatas = datasAndEnvelopes.map(([data]) => data);
  const envelopes = datasAndEnvelopes.map(([_, env]) => env);

  const combinedEnvelopes = Buffer.concat(envelopes);
  const chunks = chunkCacheTestChopEnvelope(
    combinedEnvelopes, 
    envelopes.length + (envelopes.length / 5),
  );

  const cache = new ChunkCache();
  let cntr = 0;
  const acc = [];
  for (const c of chunks) {
    cntr++;
    cache.push(c);
    if (cntr === 3) {
      acc.push(cache.nextMessageBuffer());
      cntr = 0;
    }
  }

  const missingCount = originDatas.length - acc.length;
  for (let i = 0; i < missingCount; i++) {
    acc.push(cache.nextMessageBuffer());
  }

  for (let i = 0; i < acc.length; i++) {
    const msgBuff = acc[i];
    const od = originDatas[i];
    expect(od.equals(msgBuff)).toBeTruthy();
  }
});

test('MessageReader ...', () => {
  
});

test('MessageWriter ...', () => {
  
});
