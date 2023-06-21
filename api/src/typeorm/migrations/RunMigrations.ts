import { Logger } from '@nestjs/common';
import { exec } from 'child_process';
import * as util from 'util';

export default async function startMigration() {
  const childProcess = util.promisify(exec);
  const logger = new Logger('Starting Migration...');
  const environment = process.env['CK_ENVIRONMENT'];
  const executionScript = (env) => `npm run typeorm${env} migration:run`;
  let command: string | null = null;
  try {
    if (environment === 'local') {
      command = executionScript('');
    } else {
      command = executionScript(environment);
    }
    if (command) {
      logger.log('migration command: ', {
        command,
      });
      const { stdout, stderr } = await childProcess(command);
      logger.log(`stdout: ${stdout}`);
      logger.log(`stderr: ${stderr}`);
    } else {
      logger.error('No valid command created', {
        environment,
      });
    }
  } catch (err) {
    logger.error('There was an error while performing migrations', err);
  }
}
