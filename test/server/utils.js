function Writer({reader, stream}) {
  return (msg) => {
    if (!stream.write(msg, 'utf8')) {
      reader.pause();
      stream.once('drain', () => {
        reader.resume();
      });
    }
  };
}

const Header = {
  parse: (msgStr) => {
    if (!msgStr) return [];
    const semicolonIdx = msgStr.indexOf(';');
    const headerStr = msgStr.slice(0, semicolonIdx);
    const comps = headerStr.split(',');
    return [
      {
        id: comps[0],
        role: comps[1],
        seq: Number(comps[2]),
        isEnd: parseIsEnd(comps[3]),
        status: Number(comps[4]),
        mediaType: comps[5],
        encoding: comps[6],
      },
      msgStr.slice(semicolonIdx + 1),
    ];
  },
  stringify: ({id, role, seq, isEnd, status, mediaType, encoding}) => {
    const comps = [
      id || '',
      role || '',
      typeof seq === 'number' ? seq : 0,
      typeof isEnd === 'boolean' ? isEnd : true,
      typeof status === 'number' ? status : 0,
      mediaType || '',
      encoding || '',
      ';',
    ];
    return comps.join(',');
  },
};

function parseIsEnd(comp) {
  switch (comp.toUpperCase().slice(0, 1)) {
    case 'F':
    case '0':
    case 'N':
      return false;
    default:
      return true;
  }
}

module.exports = {
  Writer,
  Header,
};
