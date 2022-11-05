import {getUsers} from './github_terraform.js';

Promise.resolve(
  (async function () {
    await getUsers();
  })()
);
