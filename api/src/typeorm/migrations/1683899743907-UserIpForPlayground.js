module.exports = class UserIpForPlayground1683899743907 {
  async up(queryRunner) {
    await queryRunner.query(
      'ALTER TABLE `USERS` ADD COLUMN `ipv4` varchar(255) UNIQUE default null'
    );
  }

  async down(queryRunner) {}
};
