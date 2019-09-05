function log(cmd, id, payload) {
  const msg = {id, role: 'log', cmd, payload};
  const msgJson = JSON.stringify(msg);
  process.stderr.write(`${msgJson}\n`, 'utf8');
}

module.exports = {
  log,
  inform: (id, payload) => log('inform', id, payload),
  fail: (id, payload) => log('fail', id, payload),
  exit: (id, payload) => log('exit', id, payload),
};