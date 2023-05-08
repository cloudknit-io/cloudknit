module.exports = class ComponentToReconciliationSkipped1683209565116 {
  async up(queryRunner) {
    await queryRunner.query(
      'ALTER TABLE `component_reconcile` ' +
        'ADD COLUMN `isSkipped` boolean DEFAULT NULL;'
    );
  }

  async down(queryRunner) {}
};
