module.exports = class EnvironmentReconcileGitSha1674766396539 {
  name = 'EnvironmentReconcileGitSha1674766396539';

  async up(queryRunner) {
    await queryRunner.query(
      `ALTER TABLE \`environment_reconcile\` ADD \`gitSha\` varchar(255) NOT NULL`
    );
  }

  async down(queryRunner) {
    await queryRunner.query(
      `ALTER TABLE \`environment_reconcile\` DROP COLUMN \`gitSha\``
    );
  }
};
