/* eslint-disable */
import * as path from 'path';
import * as fs from 'fs';
process.env.NODE_ENV = process.env.NODE_ENV || 'local';

let conf = path.resolve(__dirname, `../../.env.${process.env.NODE_ENV}`);
if(fs.existsSync(conf))
  require('dotenv').config({ path: conf });
