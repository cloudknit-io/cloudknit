module.exports = class TestMigration1687368605465 {
    async up(queryRunner) {
        await queryRunner.query(
          'SELECT * FROM `component_reconcile`;'
        );
      }
    
      async down(queryRunner) {}
}
