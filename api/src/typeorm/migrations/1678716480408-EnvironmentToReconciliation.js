module.exports = class EnvironmentToReconciliation1678716480410 {

    async up(queryRunner) {
        await queryRunner.query(
            'ALTER TABLE `environment` ' +
                'DROP COLUMN `errorMessage`,' +
                'DROP COLUMN `last_reconcile_datetime`,' +
                'DROP COLUMN `duration`,' +
                'DROP COLUMN `status`,' +
                'DROP COLUMN `estimated_cost`,' +
                'DROP COLUMN `dag`;'
          );
          await queryRunner.query(
            'ALTER TABLE `environment_reconcile` '+
                'ADD COLUMN `errorMessage` json DEFAULT NULL,' +
                'ADD COLUMN `estimated_cost` decimal(10,3) NOT NULL DEFAULT \'0.000\',' +
                'ADD COLUMN `dag` json DEFAULT NULL;'
          );
    }

    async down(queryRunner) {
    }

}
