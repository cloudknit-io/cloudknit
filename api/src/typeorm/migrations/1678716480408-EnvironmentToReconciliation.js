module.exports = class EnvironmentToReconciliation1678716480412 {
  async up(queryRunner) {
    await queryRunner.query(
      'ALTER TABLE `environment` ' +
        'DROP COLUMN `errorMessage`,' +
        'DROP COLUMN `duration`,' +
        'DROP COLUMN `status`,' +
        'DROP COLUMN `estimated_cost`,' +
        'ADD COLUMN `latest_env_recon_id` int DEFAULT NULL;'
    );
    await queryRunner.query(
      'ALTER TABLE `environment` ADD CONSTRAINT `FK_latest_env_recon` FOREIGN KEY (`latest_env_recon_id`) REFERENCES `environment_reconcile`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION'
    );
    await queryRunner.query(
      'ALTER TABLE `environment_reconcile` ' +
        'ADD COLUMN `errorMessage` json DEFAULT NULL,' +
        "ADD COLUMN `estimated_cost` decimal(10,3) NOT NULL DEFAULT '0.000'," +
        'ADD COLUMN `dag` json DEFAULT NULL;'
    );
  }

  async down(queryRunner) {}
};
