module.exports = class TeardownProtection1688573274722 {
  async up(queryRunner) {
    await queryRunner.query(
      'ALTER TABLE `team` ' +
        'ADD COLUMN `teardownProtection` boolean DEFAULT false;'
    );
  }

  async down(queryRunner) {}
};
