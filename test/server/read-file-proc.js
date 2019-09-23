const path = require('path');
const readFileB64 = require('./read-file-b64');
const getMediaType = require('./get-media-type');
const IO = require('../../node-approx/io');
const io = new IO();

const basePath = process.env.BASE_PATH;

io.read(async (message) => {
  const req = JSON.parse(message.data);
  let filePath = req.url.path;
  if (filePath === '/') {
    filePath = '/index.html';
  }

  const extension = path.extname(filePath);
  const mediaType = await getMediaType(extension);

  try {
    const data = await readFileB64(basePath, filePath);
    io.logger.info(`Sending ${data.length} bytes base64 data`);
    io.send({
      id: message.id,
      role: 'file-content',
      status: 200,
      mediaType,
      encoding: 'base64',
      data,
    });
  } catch(err) {
    if (err.code === 'ENOENT') {
      io.send({
        id: message.id, 
        role: 'file-content',
        status: 404,
      });
    } else {
      const errObj = {error: err.toString()};
      const data = JSON.stringify(errObj);
      io.send({
        id: message.id, 
        role: 'file-content',
        status: 500,
        mediaType: 'application/json',
        encoding: 'utf8',
        data,
      });
    }
  }
});
