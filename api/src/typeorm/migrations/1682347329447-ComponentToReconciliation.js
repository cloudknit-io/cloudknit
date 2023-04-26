module.exports = class ComponentToReconciliation1682347329447 {
  async up(queryRunner) {
    await queryRunner.query(
      'ALTER TABLE `components` ' +
        'DROP COLUMN `lastWorkflowRunId`,' +
        'DROP COLUMN `isDestroyed`,' +
        'DROP COLUMN `cost_resources`,' +
        'DROP COLUMN `duration`,' +
        'DROP COLUMN `status`,' +
        'DROP COLUMN `estimated_cost`,' +
        'ADD COLUMN `is_deleted` boolean DEFAULT false' +
        'ADD COLUMN `latest_comp_recon_id` int DEFAULT NULL;'
    );
    await queryRunner.query(
      'ALTER TABLE `environment` DROP CONSTRAINT `FK_latest_env_recon`;'
    );
    await queryRunner.query(
      'ALTER TABLE `component_reconcile` ' +
        'MODIFY `startDateTime` datetime DEFAULT null' +
        'ADD COLUMN `lastWorkflowRunId` varchar(255) DEFAULT NULL,' +
        'ADD COLUMN `cost_resources` json DEFAULT NULL,' +
        "ADD COLUMN `estimated_cost` decimal(10,3) NOT NULL DEFAULT '0.000'," +
        'ADD COLUMN `isDestroyed` boolean DEFAULT NULL;'
    );
  }

  async down(queryRunner) {}
};
