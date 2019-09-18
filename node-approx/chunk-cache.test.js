import ChunkCache from './chunk-cache';
import { 
  chunkCacheTestEnvelope,
  chunkCacheTestChopEnvelope,
  letterNames,
} from './test-utils';

test('ChunkCache simple test', () => {
  const data = Buffer.from('TEST-DATA');
  const lengthBuff = Buffer.from(`${data.length}:`);
  const envelope = Buffer.concat([lengthBuff, data]);
  const cache = new ChunkCache();
  cache.push(envelope);
  const msgBuff = cache.nextMessageBuffer();
  expect(data.equals(msgBuff)).toBeTruthy();
});

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