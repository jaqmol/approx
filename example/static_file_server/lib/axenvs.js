const axmsg = require('./axmsg');
const fs = require('fs');

class Envs {
  constructor(processorName, requiredEnvs=[], optionalEnvs=[]) {
    let [required, missing] = readAllEnvs(requiredEnvs);
	  logFatalIfNeeded(processorName, missing);
	  const [optional] = readAllEnvs(optionalEnvs);
	  let [ins, missingIns] = readIndexedEnvs("IN");
	  ins = amendIfNeeded(ins, missingIns, "stdin");
	  let [outs, missingOuts] = readIndexedEnvs("OUT");
	  outs = amendIfNeeded(outs, missingOuts, "stdout");
	   missing = missingIns.concat(missingOuts);
    logFatalIfNeeded(processorName, missing);

    this._processorName = processorName;
    this.required = required;
    this.optional = optional;
    this.ins = ins;
    this.outs = outs;
  }

  insOuts() {
		return [inputs(this.ins), outputs(this.outs)];
  }
}

function inputs(inPipePaths) {
	const ins = [];
	for (let inPath of inPipePaths) {
		let inStream;
		if (inPath === 'stdin') {
			inStream = process.stdin;
		} else { 
			inStream = fs.createReadStream(
				inPath, {encoding: 'utf8'},
			);
		}
		ins.push(new axmsg.Reader(inStream));
	}
	return ins;
}

function outputs(outPipePaths) {
	const outs = [];
	for (let outPath of outPipePaths) {
		let outStream;
		if (outPath === 'stdout') {
			outStream = process.stdin;
		} else { 
			outStream = fs.createWriteStream(
				outPath, {encoding: 'utf8'},
			);
		}
		outs.push(new axmsg.Writer(outStream));
	}
	return outs;
}

function logFatalIfNeeded(processorName, missing) {
	if (missing.length) {
    const errors = new axmsg.Errors(processorName);
    const err = new Error(`Required envs ${missing.join(", ")} not found`);
		errors.logFatal(null, err)
	}
}

function readAllEnvs(envNames) {
  const results = {}
	const missing = []
  for (let name of envNames) {
    const value = process.env[name];
		if (value && value.length) {
			results[name] = value;
		} else {
			missing.push(name);
		}
	}
	return [results, missing];
}

function readIndexedEnvs(prefix) {
	const results = [];
	const missing = [];

	const countName = prefix + "_COUNT";
	const countStr = process.env[countName];
	if (!countStr || !countStr.length) {
		return [results, missing];
	}

	const count = Number.parseInt(countStr);
	for (var i = 0; i < count; i++) {
		const name = `${prefix}_${i}`;
		const value = process.env[countName];
		if (value && value.length) {
			results.push(value);
		} else {
			missing.push(name);
		}
	}
	return [results, missing];
}

function amendIfNeeded(values, missing, amendment) {
	if (values.length == 0 && missing.length == 0) {
    values.push(amendment);
	}
	return values;
}

module.exports = {
  Envs,
};