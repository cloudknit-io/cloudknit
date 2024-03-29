import { Logger } from '@nestjs/common';
import { exec } from 'child_process';
import * as util from 'util';

export default async function startMigration() {
  const childProcess = util.promisify(exec);
  const logger = new Logger('Starting Migration...');
  const environment = process.env['CK_ENVIRONMENT'];
  let command: string = "npm run typeorm:dynamic migration:run";
  try {
      logger.log('migration command: ', {
        command,
      });
      const { stdout, stderr } = await childProcess(command);
      logger.log(`stdout: ${stdout}`);
      logger.log(`stderr: ${stderr}`);
  } catch (err) {
    logger.error('There was an error while performing migrations', err);
  }
}
