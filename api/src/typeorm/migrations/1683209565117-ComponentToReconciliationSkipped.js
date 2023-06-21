
module.exports = class ComponentToReconciliationSkipped1687278277726 {
  async up(queryRunner) {
    await queryRunner.query(
      'SELECT * FROM `component_reconcile`'
    );
  }

  async down(queryRunner) {}
};
